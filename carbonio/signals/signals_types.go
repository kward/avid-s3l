package signals

//go:generate stringer -output=signals_string.go -type=DirectionEnum,ConnectorEnum,FormatEnum,LevelEnum signals_types.go

// DirectionEnum of the signal.
type DirectionEnum int

const (
	unknownDirection DirectionEnum = iota
	Input
	Output
)

// ConnectorEnum physical type.
type ConnectorEnum int

const (
	unknownConnector ConnectorEnum = iota
	XLR
	Jack
)

// FormatEnum of the signal.
type FormatEnum int

const (
	unknownFormat FormatEnum = iota
	Analog
	AES
)

// LevelEnum of the signal.
type LevelEnum int

const (
	unknownLevel LevelEnum = iota
	Line
	Mic
)
