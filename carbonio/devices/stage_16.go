package devices

import (
	"fmt"
	"net"

	"github.com/kward/avid-s3l/carbonio/leds"
	"github.com/kward/avid-s3l/carbonio/signals"
)

const (
	stage16_numMicInputs   = 16
	stage16_numLineOutputs = 8
	stage16_numAESOutputs  = 4 // Each physical connector supports two signals.
)

type Stage16 struct {
	opts *options

	leds      *leds.LEDs
	micInputs signals.Signals
}

// Verify that the interface is implemented properly.
var _ Device = new(Stage16)

// NewStage16 returns a populated Stage16 struct.
func NewStage16(opts ...func(*options) error) (*Stage16, error) {
	o := &options{}
	if err := setDeviceOptions(o); err != nil {
		return nil, err
	}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	if err := o.validate(); err != nil {
		return nil, err
	}

	d := &Stage16{opts: o}

	l, err := leds.New(
		leds.SPIDelayRead(o.spiDelayRead),
		leds.SPIBaseDir(o.spiBaseDir),
		leds.Verbose(o.verbose),
	)
	if err != nil {
		return nil, err
	}
	d.leds = l

	s, err := signals.MicInputs(stage16_numMicInputs,
		signals.SPIDelayRead(o.spiDelayRead),
		signals.SPIBaseDir(o.spiBaseDir),
		signals.Verbose(o.verbose),
	)
	if err != nil {
		return nil, err
	}
	d.micInputs = s

	return d, nil
}

// LEDs implements Device.
func (d *Stage16) LEDs() *leds.LEDs { return d.leds }

// NumMicInputs implements Device.
func (d *Stage16) NumMicInputs() int { return stage16_numMicInputs }

// MicInput returns the signal for the specified input number.
func (d *Stage16) MicInput(input int) (*signals.Signal, error) {
	if input < 1 || input > stage16_numMicInputs {
		return nil, fmt.Errorf("invalid input number %d", input)
	}
	return d.micInputs[input], nil
}

func (d *Stage16) IP() net.IP { return d.opts.ip }
