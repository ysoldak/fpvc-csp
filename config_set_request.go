package csp

type ConfigSetRequest struct {
	ID     byte
	Offset byte
	Data   []byte
}

func NewConfigSetRequest(id, offset byte, data []byte) *ConfigSetRequest {
	return &ConfigSetRequest{
		ID:     id,
		Offset: offset,
		Data:   data,
	}
}

func NewConfigSetRequestFromMessage(message *Message) *ConfigSetRequest {
	return &ConfigSetRequest{
		ID:     message.Payload[0],
		Offset: message.Payload[1],
		Data:   message.Payload[2:],
	}
}

func (c *ConfigSetRequest) Message() *Message {
	data := make([]byte, 2+len(c.Data))
	data[0] = c.ID
	data[1] = c.Offset
	copy(data[2:], c.Data)
	return NewRequest(CommandConfigSet, data)
}
