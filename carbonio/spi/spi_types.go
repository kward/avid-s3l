package spi

//go:generate stringer -output=spi_string.go -type=Enum spi_types.go

type Enum int

const (
	unknownSPI Enum = iota

	// LEDs.
	Power
	Status
	Mute

	// Inputs.
	Gain
	Pad
	Phantom

	// TestLEDs (for testing only).
	Blinky
)
