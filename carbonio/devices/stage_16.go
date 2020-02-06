package devices

import (
	"fmt"

	"github.com/kward/avid-s3l/carbonio/leds"
	"github.com/kward/avid-s3l/carbonio/signals"
)

const (
	stage16_micInputs   = 16
	stage16_lineOutputs = 8
	stage16_aesOutputs  = 2
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

	s, err := signals.MicInputs(stage16_micInputs)
	if err != nil {
		return nil, err
	}
	d.micInputs = s

	return d, nil
}

// NumMicInputs implements Device.
func (d *Stage16) NumMicInputs() uint { return stage16_micInputs }

func (d *Stage16) MicInput(input uint) (*signals.Signal, error) {
	if input < 1 || input > stage16_micInputs {
		return nil, fmt.Errorf("invalid input number %d", input)
	}
	return d.micInputs[input], nil
}
