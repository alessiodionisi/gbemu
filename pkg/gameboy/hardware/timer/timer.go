package timer

import "github.com/adnsio/gbemu/pkg/gameboy/bits"

const (
	ControlClockSelect0 = 0
	ControlClockSelect1 = 1
	ControlEnabled      = 2
)

type Timer struct {
	Control         uint8
	Counter         uint8
	Modulo          uint8
	DividerRegister uint8
	InternalCounter int
}

func NewTimer() *Timer {
	return &Timer{}
}

func (t *Timer) Update(cycles int) {
	// todo divider register

	isEnabled := bits.Test(t.Control, ControlEnabled)

	if isEnabled {
		t.InternalCounter -= cycles

		if t.InternalCounter <= 0 {
			frequency := t.Counter & 0x3

			switch frequency {
			case 0:
				t.InternalCounter = 1024
			case 1:
				t.InternalCounter = 16
			case 2:
				t.InternalCounter = 64
			case 3:
				t.InternalCounter = 256
			}

			if t.Counter == 255 {
				t.Counter = t.Modulo
				// todo request interrupt 2
			} else {
				t.Counter++
			}
		}
	}
}
