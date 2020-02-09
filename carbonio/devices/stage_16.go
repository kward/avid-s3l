package devices

import (
	"fmt"

	"github.com/kward/avid-s3l/carbonio/leds"
	"github.com/kward/avid-s3l/carbonio/signals"
)

const (
	stage16_numMicInputs   = 16
	stage16_numLineOutputs = 8
	stage16_numAESOutputs  = 2
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
	o.setSPIBaseDir(SPIDevicesDir)
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	d := &Stage16{
		opts:      o,
		powerLED:  leds.Power,
		statusLED: leds.Status,
		muteLED:   leds.Mute,
	}

	s, err := signals.MicInputs(o.spiBaseDir, stage16_numMicInputs)
	if err != nil {
		return nil, err
	}
	d.micInputs = s

	return d, nil
}

// NumMicInputs implements Device.
func (d *Stage16) NumMicInputs() uint { return stage16_numMicInputs }

func (d *Stage16) MicInput(input uint) (*signals.Signal, error) {
	if input < 1 || input > stage16_numMicInputs {
		return nil, fmt.Errorf("invalid input number %d", input)
	}
	return d.micInputs[input], nil
}
