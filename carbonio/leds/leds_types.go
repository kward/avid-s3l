package leds

//go:generate stringer -output leds_string.go -type=LEDState leds_types.go

type LEDState int

const (
	Unknown LEDState = iota
	Off
	Alert
	On
	testLEDState LEDState = 255
)
