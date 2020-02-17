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
	"path"

	"github.com/kward/avid-s3l/carbonio/helpers"
)

type value byte
type byState map[State]value
type byValue map[value]State

// LED describes a Carbon I/O LED.
type LED struct {
	opts   *options
	name   string
	iface  string
	spi    string
	states byState
	values byValue
}

type LEDs struct {
	power, status, mute *LED
}

func (l LEDs) Power() *LED  { return l.power }
func (l LEDs) Status() *LED { return l.status }
func (l LEDs) Mute() *LED   { return l.mute }

// New instantiates the
func New(opts ...func(*options) error) (*LEDs, error) {
	leds := &LEDs{new(LED), new(LED), new(LED)}

	for _, led := range []struct {
		led    **LED
		name   string
		iface  string
		states byState
	}{
		{&leds.power, "Power", "spi4.0/status_led_1_en", byState{Off: '0', Alert: '1', On: '2', testState: 255}},
		{&leds.status, "Status", "spi4.0/status_led_0_en", byState{Off: '0', Alert: '1', On: '2', testState: 255}},
		{&leds.mute, "Mute", "spi4.0/mute_led_en", byState{Off: '0', On: '1', testState: 255}},
	} {
		l, err := NewLED(led.name, led.iface, led.states, opts...)
		if err != nil {
			return nil, err
		}
		*led.led = l
	}

	return leds, nil
}

// NewLED instantiates a new LED.
func NewLED(name string, iface string, states byState, opts ...func(*options) error) (*LED, error) {
	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	if err := o.validate(); err != nil {
		return nil, err
	}

	led := &LED{
		name:   name,
		iface:  iface,
		spi:    path.Join(o.spiBaseDir, iface),
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
	if l == nil {
		return Unknown, fmt.Errorf("LED uninitialized")
	}

	if l.opts.verbose {
		fmt.Printf("reading LED from %s\n", l.spi)
	}
	v, err := helpers.ReadByte(l.spi)
	if err != nil {
		return Unknown, err
	}
	s, ok := l.values[value(v)]
	if !ok {
		return Unknown, fmt.Errorf("unrecognized LED value %q [%d]", v, v)
	}
	return s, nil
}

// SetState changes the state of the LED.
func (l *LED) SetState(s State) error {
	if l == nil {
		return fmt.Errorf("LED uninitialized")
	}

	var v value
	v, ok := l.states[s]
	if !ok {
		return fmt.Errorf("unrecognized LED state %q [%d]", s, s)
	}

	return helpers.WriteByte(l.spi, byte(v))
}

func (l *LED) Path() string {
	if l == nil {
		return "uninitialized"
	}

	return l.iface
}

// String provides a human readable state output.
func (l *LED) String() string {
	if l == nil {
		return "uninitialized"
	}

	s, _ := l.State()
	return fmt.Sprintf("%s", s)
}
