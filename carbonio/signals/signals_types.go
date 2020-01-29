// Package signals represents different types of physical signals available.
package signals

//go:generate stringer -output=signal_string.go -type=Connector,Format,Level

type Connector int

const (
	XLRConnector Connector = iota
	JackConnector
)

type Format int

const (
	AnalogFormat Format = iota
	AESFormat
)

type Level int

const (
	LineLevel Level = iota
	MicLevel
)

type Pad int

const (
	PadDisabled = iota
	PadEnabled
)
