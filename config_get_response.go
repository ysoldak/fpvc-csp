package csp

type ConfigGetResponse struct {
	ID     byte
	Offset byte
	Data   []byte
}

func NewConfigGetResponse(id, offset byte, data []byte) *ConfigGetResponse {
	return &ConfigGetResponse{
		ID:     id,
		Offset: offset,
		Data:   data,
	}
}

func NewConfigGetResponseFromMessage(message *Message) *ConfigGetResponse {
	return &ConfigGetResponse{
		ID:     message.Payload[0],
		Offset: message.Payload[1],
		Data:   message.Payload[2:],
	}
}

func (c *ConfigGetResponse) Message() *Message {
	data := make([]byte, 2+len(c.Data))
	data[0] = c.ID
	data[1] = c.Offset
	copy(data[2:], c.Data)
	return NewResponse(CommandConfigGet, data)
}
