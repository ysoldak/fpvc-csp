package csp

type Beacon struct {
	ID          byte
	Name        string
	Description string
}

func NewBeacon(id byte, name, description string) *Beacon {
	return &Beacon{
		ID:          id,
		Name:        name,
		Description: description,
	}
}

func NewBeaconFromMessage(message *Message) *Beacon {
	return &Beacon{
		ID:          message.Payload[0],
		Name:        string(message.Payload[1:11]),
		Description: string(message.Payload[11:31]),
	}
}

func (b *Beacon) Message() *Message {
	data := make([]byte, 1+10+20)
	data[0] = b.ID
	nameLen := len(b.Name)
	for i := 0; i < 10; i++ {
		if i < nameLen {
			data[1+i] = b.Name[i]
		} else {
			data[1+i] = 0x20
		}
	}
	descLen := len(b.Description)
	for i := 0; i < 20; i++ {
		if i < descLen {
			data[11+i] = b.Description[i]
		} else {
			data[11+i] = 0x20
		}
	}

	return NewBroadcast(CmdBeacon, data)
}
