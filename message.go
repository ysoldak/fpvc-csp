package csp

// Commands
type Command byte

const (
	// < ID[1], Name[10], Description[20] <-- when not in battle, broadcast beacon message
	CommandBeacon Command = 0x60

	// > ID[1], Offset[1], Length[1]   <-- Length[1] = 1 to 110; to read full config: Offset[1] = 0, Length[1] = 110
	// < ID[1], Offset[1], Data[<=110] <-- data length = message payload length - 2, shall be equal to Length[1] passed in request
	CommandConfigGet Command = 0x72

	// > ID[1], Offset[1], Data[<=110] <-- data length = message payload length - 2
	// < ID[1], Offset[1], Data[<=110] <-- ID[1] can be new if team or player changed
	CommandConfigSet Command = 0x74

	// > ID[1], Lives[1] <-- the player who was hit sends this message
	// < ID[1], Power[1] <-- the player who was shooting replies with this message
	CommandHit Command = 0x82
)

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

func (m *Message) Bytes() []byte {
	b := make([]byte, 5+len(m.Payload)+1)
	b[0] = m.Header[0]
	b[1] = m.Header[1]
	b[2] = byte(m.Direction)
	b[3] = m.Length
	b[4] = byte(m.Command)
	copy(b[5:], m.Payload)
	b[len(b)-1] = m.Checksum
	return b
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
