package csp

type ConfigSetResponse struct {
	ID     byte
	Offset byte
	Data   []byte
}

func NewConfigSetResponse(id, offset byte, data []byte) *ConfigSetResponse {
	return &ConfigSetResponse{
		ID:     id,
		Offset: offset,
		Data:   data,
	}
}

func NewConfigSetResponseFromMessage(message *Message) *ConfigSetResponse {
	return &ConfigSetResponse{
		ID:     message.Payload[0],
		Offset: message.Payload[1],
		Data:   message.Payload[2:],
	}
}

func (c *ConfigSetResponse) Message() *Message {
	data := make([]byte, 2+len(c.Data))
	data[0] = c.ID
	data[1] = c.Offset
	copy(data[2:], c.Data)
	return NewResponse(CmdConfigSet, data)
}
