package csp

type HitRequest struct {
	ID    byte
	Lives byte
}

func NewHitRequest(id, lives byte) *HitRequest {
	return &HitRequest{
		ID:    id,
		Lives: lives,
	}
}

func NewHitRequestFromMessage(message *Message) *HitRequest {
	return &HitRequest{
		ID:    message.Payload[0],
		Lives: message.Payload[1],
	}
}

func (h *HitRequest) Message() *Message {
	return NewRequest(CmdHit, []byte{h.ID, h.Lives})
}
