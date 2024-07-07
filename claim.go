package csp

type Claim struct {
	PlayerID byte
	Power    byte
}

func NewClaim(id, power byte) Claim {
	return Claim{
		PlayerID: id,
		Power:    power,
	}
}

func NewClaimFromMessage(message Message) Claim {
	return Claim{
		PlayerID: message.Payload[0],
		Power:    message.Payload[1],
	}
}

func (c Claim) Message() Message {
	return NewMessage(COMMAND_CLAIM, []byte{c.PlayerID, c.Power})
}
