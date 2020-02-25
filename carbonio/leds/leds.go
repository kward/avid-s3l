/*
Package leds enables control of LEDs through the Carbon I/O front-panel LEDs.

The Avid Stage 16 LEDs are controlled through the kernel interface of the spi
device. The current state of the LED can be read by looking at the contents of
the file that represents the device interface. Changing the state can be done by
writing to the same interface.
*/

package leds

import (
	"fmt"
	"os"

	"github.com/kward/avid-s3l/carbonio/spi"
)

type byState map[State]int
type byValue map[int]State

type LEDs struct {
	power, status, mute *LED
}

func (l LEDs) Power() *LED  { return l.power }
func (l LEDs) Status() *LED { return l.status }
func (l LEDs) Mute() *LED   { return l.mute }

// New instantiates the package level LEDs.
func New(opts ...func(*options) error) (*LEDs, error) {
	leds := &LEDs{new(LED), new(LED), new(LED)}

	for _, led := range []struct {
		led    **LED
		enum   spi.Enum
		states byState
	}{
		{&leds.power, spi.PowerLED, byState{Off: 0, Alert: 1, On: 2, testState: 255}},
		{&leds.status, spi.StatusLED, byState{Off: 0, Alert: 1, On: 2, testState: 255}},
		{&leds.mute, spi.MuteLED, byState{Off: 0, On: 1, testState: 255}},
	} {
		l, err := NewLED(led.enum, led.states, opts...)
		if err != nil {
			return nil, err
		}
		*led.led = l
	}

	return leds, nil
}

// LED describes a Carbon I/O LED.
type LED struct {
	opts   *options
	spi    *spi.SPI
	states byState
	values byValue
}

// Ensure spi interfaces are implemented.
var _ spi.Implementation = new(LED)

// NewLED instantiates a new LED.
func NewLED(enum spi.Enum, states byState, opts ...func(*options) error) (*LED, error) {
	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	if err := o.validate(); err != nil {
		return nil, err
	}

	s, err := spi.New(enum, 0,
		spi.DelayRead(o.spiDelayRead),
		spi.BaseDir(o.spiBaseDir),
	)
	if err != nil {
		return nil, fmt.Errorf("failure instantiating %s SPI; %s", enum.String(), err)
	}
	led := &LED{
		spi:    s,
		states: states,
		values: byValue{},
		opts:   o,
	}
	for k, v := range led.states {
		led.values[v] = k
	}
	if o.verbose {
		fmt.Fprintf(os.Stderr, "%#v\n", led)
	}
	return led, nil
}

// State returns the active state of the LED.
func (l *LED) State() (State, error) {
	if l.opts.verbose {
		fmt.Printf("reading %s LED from %s\n", l.spi.Name(), l.spi.Path())
	}
	v, err := l.spi.Read()
	if err != nil {
		return Unknown, err
	}
	s, ok := l.values[v]
	if !ok {
		return Unknown, fmt.Errorf("unrecognized %s LED value %d", l.spi.Name(), v)
	}
	return s, nil
}

// SetState changes the state of the LED.
func (l *LED) SetState(s State) error {
	var v int
	v, ok := l.states[s]
	if !ok {
		return fmt.Errorf("unrecognized %s LED state %q [%d]", l.spi.Name(), s, s)
	}

	return l.spi.Write(v)
}

// Initialize implements spi.Implementation.
func (l *LED) Initialize() error { return l.SetState(Off) }

// Name implements spi.Implementation.
func (l *LED) Name() string { return l.spi.Name() }

// Path implements spi.Implementation.
func (l *LED) Path() string { return l.spi.Path() }

// Raw implements spi.Implementation.
func (l *LED) Raw() []byte { return l.spi.Raw() }

// String provides a human readable state output.
func (l *LED) String() string {
	s, _ := l.State()
	return fmt.Sprintf("%s", s)
}
