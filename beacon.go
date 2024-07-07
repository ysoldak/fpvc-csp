package csp

type Beacon struct {
	PlayerID    byte
	Name        string
	Description string
}

func NewBeacon(id byte, name, description string) Beacon {
	return Beacon{
		PlayerID:    id,
		Name:        name,
		Description: description,
	}
}

func NewBeaconFromMessage(message Message) Beacon {
	return Beacon{
		PlayerID:    message.Payload[0],
		Name:        string(message.Payload[1:11]),
		Description: string(message.Payload[11:31]),
	}
}

func (b Beacon) Message() Message {
	data := make([]byte, 1+10+20)
	data[0] = b.PlayerID
	copy(data[1:], []byte(b.Name))
	copy(data[11:], []byte(b.Description))
	return NewMessage(COMMAND_BEACON, data)
}
