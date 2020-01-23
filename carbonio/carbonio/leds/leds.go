// Package leds provides definitions for the CarbonIO front-panel LEDs.
package leds

//go:generate stringer -type=LEDState

const (
	PowerIface  = "/sys/bus/spi/devices/spi4.0/status_led_1_en"
	StatusIface = "/sys/bus/spi/devices/spi4.0/status_led_0_en"
	MuteIface   = "/sys/bus/spi/devices/spi4.0/mute_led_en"
)

type LEDState int

const (
	Off LEDState = iota
	On
	Alert
	Green  = On
	Orange = Alert
)

type LED interface {
	State() (*LEDState, error)
	SetState() error
}

type Power struct {
	state LEDState
}

var _ LED = new(Power)

func (l *Power) State() (*LEDState, error) {
	return nil, nil
}
func (l *Power) SetState() error {
	return nil
}
