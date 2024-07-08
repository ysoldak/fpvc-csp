package csp

type HitResponse struct {
	ID    byte
	Power byte
}

func NewHitResponse(id, power byte) *HitResponse {
	return &HitResponse{
		ID:    id,
		Power: power,
	}
}

func NewHitResponseFromMessage(message *Message) *HitResponse {
	return &HitResponse{
		ID:    message.Payload[0],
		Power: message.Payload[1],
	}
}

func (c *HitResponse) Message() *Message {
	return NewResponse(CmdHit, []byte{c.ID, c.Power})
}
