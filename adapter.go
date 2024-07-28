package csp

import (
	"errors"
	"io"
	"time"
)

var ErrNoData = errors.New("no data available")
var ErrWrongChecksum = errors.New("wrong checksum")
var ErrWrite = errors.New("write failed")
var ErrWriteLength = errors.New("write failed to send all bytes")
var ErrTimeout = errors.New("timeout")

const BeaconInterval = 6 * time.Second

const maxPayload = 1 + 1 + 110 // CONFIG SET requests: ID[1], Offset[1], Data[up to 110 bytes]

const (
	stateIdle byte = iota
	stateHeader
	stateDirection
	stateLength
	stateCommand
	statePayload
	stateChecksum
)

type Adapter struct {
	wire io.ReadWriter

	lowestID            byte
	beaconReferenceTime int64

	state   byte
	message Message
}

func NewAdapter(wire io.ReadWriter) *Adapter {
	return &Adapter{
		wire: wire,
		message: Message{
			Payload: []byte{},
		},
	}
}

// Send a message.
func (a *Adapter) Send(message *Message) error {
	bytes := make([]byte, 5+maxPayload+1) // optimisation to avoid heap allocation: could allocate only required size, but that is not constant
	_ = message.Bytes(bytes)
	logTs("SEND ")
	for _, b := range bytes {
		log(" %02X", b)
	}
	log("\n")
	n, err := a.wire.Write(bytes[0:message.Size()])
	if err != nil {
		return ErrWrite
	}
	if n != len(bytes) {
		return ErrWriteLength
	}
	return nil
}

// Receive a message; returns nil if no message is available (yet).
func (a *Adapter) Receive(result *Message) error {
	buf := make([]byte, 16)
	for {
		n, err := a.wire.Read(buf)
		if err != nil || n == 0 {
			return ErrNoData
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
					a.state = stateDirection
				} else {
					a.state = stateIdle
				}
			case stateDirection:
				logTs("DIRECTION %02X\n", b)
				if b != byte(DirRequest) && b != byte(DirResponse) {
					a.state = stateIdle
					continue
				}
				a.message.Direction = Direction(b)
				a.state = stateLength
			case stateLength:
				logTs("LENGTH %02X\n", b)
				if b > maxPayload {
					a.state = stateIdle
					continue
				}
				a.message.Length = b
				a.message.Payload = a.message.Payload[:0]
				a.message.Checksum = b
				a.state = stateCommand
			case stateCommand:
				logTs("COMMAND %02X\n", b)
				a.message.Command = Command(b)
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
				for _, bb := range a.message.Payload {
					log(" %02X", bb)
				}
				log("\n")
				logTs("CHECKSUM expected %02X ?= %02X actual\n", a.message.Checksum, b)
				result.Copy(&a.message)
				a.state = stateIdle
				if result.Checksum == b {
					a.handleBeaconMaybe(result)
					return nil
				} else {
					return ErrWrongChecksum
				}
			}
		}
	}
}

// Reset the state machine and clear the message buffer.
func (a *Adapter) Reset() {
	a.state = stateIdle
	buf := make([]byte, 16)
	for {
		n, err := a.wire.Read(buf)
		if err != nil || n == 0 {
			return
		}
	}
}
