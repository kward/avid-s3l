package spi

//go:generate stringer -output=spi_string.go -type=Enum spi_types.go

type Enum int

const (
	unknownSPI Enum = iota

	// LEDs.
	PowerLED
	StatusLED
	MuteLED

	// ADC / Inputs.
	Gain // 0 if uninitialized.
	Pad
	Phantom

	// DAC / Outputs.
	Attenuation // 255 if uninitialized.
	Mute
	OpAmp
	Phase

	// Other.
	Switch

	// TestLEDs (for testing only).
	Blinky
)
