package csp

type Hit struct {
	PlayerID byte
	Lives    byte
}

func NewHit(id, lives byte) Hit {
	return Hit{
		PlayerID: id,
		Lives:    lives,
	}
}

func NewHitFromMessage(message Message) Hit {
	return Hit{
		PlayerID: message.Payload[0],
		Lives:    message.Payload[1],
	}
}

func (h Hit) Message() Message {
	return NewMessage(COMMAND_HIT, []byte{h.PlayerID, h.Lives})
}
