package leds

//go:generate stringer -type=LEDState ledstates.go

type LEDState int

const (
	Unknown LEDState = iota
	Off
	Alert
	On
)
