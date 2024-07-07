package csp

import (
	"errors"
	"io"
)

var ErrNoData = errors.New("no data available")
var ErrWrongChecksum = errors.New("wrong checksum")
var ErrWrite = errors.New("write failed")
var ErrWriteLength = errors.New("write failed to send all bytes")

const (
	stateIdle byte = iota
	stateHeader
	stateLength
	stateCommand
	statePayload
	stateChecksum
)

type Adapter struct {
	wire io.ReadWriter

	state   byte
	message Message
}

func NewAdapter(wire io.ReadWriter) *Adapter {
	return &Adapter{
		wire: wire,
	}
}

// Send a message.
func (a *Adapter) Send(m *Message) error {
	bytes := m.Bytes()
	logTs("SEND ")
	for _, b := range bytes {
		log(" %02X", b)
	}
	log("\n")
	n, err := a.wire.Write(bytes)
	if err != nil {
		return ErrWrite
	}
	if n != len(bytes) {
		return ErrWriteLength
	}
	return nil
}

// Receive a message; returns nil if no message is available (yet).
func (a *Adapter) Receive() (*Message, error) {
	buf := make([]byte, 16)
	for {
		n, err := a.wire.Read(buf)
		if err != nil || n == 0 {
			return nil, ErrNoData
		}
		for i := 0; i < n; i++ {
			b := buf[i]
			switch a.state {
			case stateIdle:
				if b == '$' {
					logTs("IDLE %02X\n", b)
					a.message.Header[0] = b
					a.state = stateHeader
				}
			case stateHeader:
				if b == 'C' {
					logTs("HEADER %02X\n", b)
					a.message.Header[1] = b
					a.state = stateLength
				} else {
					a.state = stateIdle
				}
			case stateLength:
				logTs("LENGTH %02X\n", b)
				if b > MAX_PAYLOAD {
					a.state = stateIdle
					continue
				}
				a.message.Length = b
				a.message.Payload = []byte{}
				a.message.Checksum = b
				a.state = stateCommand
			case stateCommand:
				logTs("COMMAND %02X\n", b)
				a.message.Command = b
				a.message.Checksum ^= b
				a.state = statePayload
			case statePayload:
				a.message.Payload = append(a.message.Payload, b)
				a.message.Checksum ^= b
				if len(a.message.Payload) == int(a.message.Length) {
					a.state = stateChecksum
				}
			case stateChecksum:
				logTs("PAYLOAD ")
				for _, bb := range a.message.Bytes() {
					log(" %02X", bb)
				}
				log("\n")
				logTs("CHECKSUM expected %02X ?= %02X actual\n", a.message.Checksum, b)
				result := a.message
				a.message = Message{}
				a.state = stateIdle
				if result.Checksum == b {
					return &result, nil
				} else {
					return nil, ErrWrongChecksum
				}
			}
		}
	}
}

// Reset the state machine and clear the message buffer.
func (a *Adapter) Reset() {
	a.state = stateIdle
	a.message = Message{}
	buf := make([]byte, 16)
	for {
		n, err := a.wire.Read(buf)
		if err != nil || n == 0 {
			return
		}
	}
}
