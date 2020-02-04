package devices

import (
	"github.com/kward/avid-s3l/carbonio/leds"
	"github.com/kward/avid-s3l/carbonio/signals"
)

const (
	numMicInputs   = 16
	numLineOutputs = 8
	numAESOutputs  = 2
)

type Stage16 struct {
	opts *options

	powerLED, statusLED, muteLED *leds.LED
	micInputs                    signals.Signals
}

// Verify that the interface is implemented properly.
var _ Device = new(Stage16)

// NewStage16 returns a populated Stage16 struct.
func NewStage16(opts ...func(*options) error) (*Stage16, error) {
	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	d := &Stage16{
		opts: o,
	}

	// LEDs
	d.powerLED = leds.Power
	d.statusLED = leds.Status
	d.muteLED = leds.Mute

	s, err := signals.MicInputs(numMicInputs)
	if err != nil {
		return nil, err
	}
	d.micInputs = s

	return d, nil
}
