//go:build avr

package csp

// Wait for a message with the given command and direction.
func (a *Adapter) Wait(command Command, direction Direction, timeout int64, message *Message) error {
	start := runtime_nanotime()
	for runtime_nanotime()-start < timeout {
		err := a.Receive(message)
		if err != nil {
			continue
		}
		// wait for correct message
		if message.Command == command && message.Direction == direction {
			return nil
		}
	}
	return ErrTimeout
}

// BeaconTime returns the next time when a beacon with the given ID should be broadcasted.
func (a *Adapter) BeaconTime(id byte) int64 {
	if a.beaconReferenceTime == 0 {
		return 0
	}
	offset := beaconOffset(id)
	t := a.beaconReferenceTime + offset
	now := runtime_nanotime()
	for t < now {
		t += int64(BeaconInterval)
	}
	return t
}

func (a *Adapter) handleBeaconMaybe(message *Message) {
	if message.Command != CmdBeacon {
		return
	}
	id := message.Payload[0]
	if a.lowestID == 0 || a.lowestID > id {
		a.lowestID = id
	}
	if id != a.lowestID {
		return
	}
	// The beacon with the lowest ID is the reference beacon.
	offset := beaconOffset(a.lowestID)
	a.beaconReferenceTime = runtime_nanotime() - offset
}

func beaconOffset(id byte) int64 {
	team := (id << 4) - 0x0A
	player := (id & 0x0F) - 1
	return int64(team)*1_000_000_000 + int64(player)*100_000_000
}
