//go:build !avr

package csp

import "time"

// Wait for a message with the given command and direction.
func (a *Adapter) Wait(command Command, direction Direction, timeout time.Duration, message *Message) error {
	start := time.Now()
	for time.Since(start) < timeout {
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
func (a *Adapter) BeaconTime(id byte) time.Time {
	if a.beaconReferenceTime == 0 {
		return time.Time{}
	}
	offset := beaconOffset(id)
	t := a.beaconReferenceTime + offset.Milliseconds()
	now := time.Now().UnixMilli()
	for t < now {
		t += BeaconInterval.Milliseconds()
	}
	return time.UnixMilli(t)
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
	a.beaconReferenceTime = time.Now().UnixMilli() - offset.Milliseconds()
}

func beaconOffset(id byte) time.Duration {
	team := (id << 4) - 0x0A
	player := (id & 0x0F) - 1
	return time.Duration(team)*time.Second + time.Duration(player)*100*time.Millisecond
}
