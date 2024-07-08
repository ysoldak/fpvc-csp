package csp

type ConfigGetRequest struct {
	ID     byte
	Offset byte
	Length byte
}

func NewConfigGetRequest(id, offset, length byte) *ConfigGetRequest {
	return &ConfigGetRequest{
		ID:     id,
		Offset: offset,
		Length: length,
	}
}

func NewConfigGetRequestFromMessage(message *Message) *ConfigGetRequest {
	return &ConfigGetRequest{
		ID:     message.Payload[0],
		Offset: message.Payload[1],
		Length: message.Payload[2],
	}
}

func (c *ConfigGetRequest) Message() *Message {
	return NewRequest(CommandConfigGet, []byte{c.ID, c.Offset, c.Length})
}
