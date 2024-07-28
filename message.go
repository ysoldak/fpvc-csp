package csp

import "errors"

// Commands
type Command byte

const (
	// < ID[1], Name[10], Description[20] <-- when not in battle, broadcast beacon message
	CmdBeacon Command = 0x60

	// > ID[1], Offset[1], Length[1]   <-- Length[1] = 1 to 110; to read full config: Offset[1] = 0, Length[1] = 110
	// < ID[1], Offset[1], Data[<=110] <-- data length = message payload length - 2, shall be equal to Length[1] passed in request
	CmdConfigGet Command = 0x72

	// > ID[1], Offset[1], Data[<=110] <-- data length = message payload length - 2
	// < ID[1], Offset[1], Data[<=110] <-- ID[1] can be new if team or player changed
	CmdConfigSet Command = 0x74

	// > ID[1], Lives[1] <-- the player who was hit sends this message
	// < ID[1], Power[1] <-- the player who was shooting replies with this message
	CmdHit Command = 0x82
)

var errBufferTooSmall = errors.New("buffer too small")

// Directions
type Direction byte

const (
	DirRequest  Direction = '>'
	DirResponse Direction = '<'
)

type Message struct {
	Header    [2]byte   // '$' + 'C'
	Direction Direction // '>' (0x3E) = Request, '<' (0x3C) = Response/Broadcast
	Length    byte      // Length of the payload
	Command   Command   // 0x82 = Claim, 0x83 = Hit, etc
	Payload   []byte    // Data
	Checksum  byte      // XOR of all bytes from length to the end of payload
}

func NewRequest(command Command, data []byte) *Message {
	return NewMessage(DirRequest, command, data)
}

func NewResponse(command Command, data []byte) *Message {
	return NewMessage(DirResponse, command, data)
}

func NewBroadcast(command Command, data []byte) *Message {
	return NewResponse(command, data)
}

func NewMessage(direction Direction, command Command, data []byte) *Message {
	checksum := byte(len(data)) ^ byte(command)
	for _, b := range data {
		checksum ^= b
	}
	return &Message{
		Header:    [2]byte{'$', 'C'},
		Direction: direction,
		Length:    byte(len(data)),
		Command:   command,
		Payload:   data,
		Checksum:  checksum,
	}
}

func (m *Message) Copy(o *Message) {
	copy(m.Header[:], o.Header[:])
	m.Direction = o.Direction
	m.Length = o.Length
	m.Command = o.Command
	if m.Payload == nil {
		m.Payload = make([]byte, o.Length)
	}
	copy(m.Payload, o.Payload)
	m.Checksum = o.Checksum
}

func (m *Message) Bytes(b []byte) error {
	if len(b) < int(m.Size()) {
		return errBufferTooSmall
	}
	b[0] = m.Header[0]
	b[1] = m.Header[1]
	b[2] = byte(m.Direction)
	b[3] = m.Length
	b[4] = byte(m.Command)
	copy(b[5:], m.Payload)
	b[m.Size()-1] = m.Checksum
	return nil
}

func (m *Message) Size() byte {
	return 5 + m.Length + 1
}

func (m *Message) IsRequest() bool {
	return m.Direction == DirRequest
}

func (m *Message) IsResponse() bool {
	return m.Direction == DirResponse
}

func (m *Message) IsBroadcast() bool {
	return m.IsResponse()
}
