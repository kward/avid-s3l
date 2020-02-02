// Package leds enables control of LEDs through the Carbon I/O front-panel LEDs.
//
// The Avid Stage 16 LEDs are controlled through the kernel interface of the
// spi device. The current state of the LED can be read by looking at the
// contents of the file that represents the device interface. Changing the
// state can be done by writing to the same interface.
package leds

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kward/golib/math"
)

const (
	maxDataLen = 40
)

type value byte
type byState map[State]value
type byValue map[value]State

// LED describes the struct for any LED.
type LED struct {
	name   string
	iface  string
	states byState
	values byValue
}

// New instantiates a new LED.
func New(name string, iface string, states byState) *LED {
	led := &LED{
		name:   name,
		iface:  iface,
		states: states,
		values: byValue{},
	}
	for k, v := range led.states {
		led.values[v] = k
	}
	return led
}

// State returns the active state of the LED.
func (l *LED) State() (State, error) {
	v, err := readValue(l.iface)
	if err != nil {
		return Unknown, err
	}
	s, ok := l.values[v]
	if !ok {
		return Unknown, fmt.Errorf("unrecognized LED value %q [%d]", v, v)
	}
	return s, nil
}

// SetState changes the state of the LED.
func (l *LED) SetState(s State) error {
	var v value
	v, ok := l.states[s]
	if !ok {
		return fmt.Errorf("unrecognized LED state %q [%d]", s, s)
	}

	return writeValue(l.iface, v)
}

// String provides a human readable state output.
func (l *LED) String() string {
	s, _ := l.State()
	return fmt.Sprintf("%s LED %s", l.name, s)
}

//=============================================================================

var readFileFn = ioutil.ReadFile

// readValue from the SPI interface.
func readValue(filename string) (value, error) {
	data, err := readFileFn(filename)
	if err != nil {
		return 0, err
	}
	if len(data) != 2 || data[1] != '\n' {
		return 0, fmt.Errorf("%q contains unexpected data; %v", filename, data[0:math.MinInt(len(data), maxDataLen)])
	}
	return value(data[0]), nil
}

var writeFileFn = ioutil.WriteFile

// writeValue to the SPI interface.
func writeValue(filename string, v value) error {
	return writeFileFn(filename, []byte{byte(v), '\n'}, os.FileMode(0644))
}

//=============================================================================

var (
	// Power provides access to the Power LED.
	Power *LED
	// Status provides access to the Status LED.
	Status *LED
	// Mute provides access to the Mute LED.
	Mute *LED
)

func init() {
	Power = New("Power", "/sys/bus/spi/devices/spi4.0/status_led_1_en",
		byState{Off: '0', Alert: '1', On: '2', testState: 255},
	)
	Status = New("Status", "/sys/bus/spi/devices/spi4.0/status_led_0_en",
		byState{Off: '0', Alert: '1', On: '2', testState: 255},
	)
	Mute = New("Mute", "/sys/bus/spi/devices/spi4.0/mute_led_en",
		byState{Off: '0', On: '1', testState: 255},
	)
}
