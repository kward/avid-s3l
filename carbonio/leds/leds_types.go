package leds

//go:generate stringer -output leds_string.go -type=State leds_types.go

type State int

const (
	Unknown State = iota
	Off
	Alert
	On
	testState State = 255
)
