package spi

//go:generate stringer -output=spi_string.go -type=Enum spi_types.go

type Enum int

const (
	unknownSPI Enum = iota

	// LEDs.
	PowerLED
	StatusLED
	MuteLED

	// Inputs.
	Gain // 0 if uninitialized.
	Pad
	Phantom

	// Outputs.
	Attenuation // 255 if uninitialized.
	Mute
	OpAmp
	Phase

	// TestLEDs (for testing only).
	Blinky
)
