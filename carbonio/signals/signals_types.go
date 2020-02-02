// Package signals represents different types of physical signals available.
package signals

//go:generate stringer -output=signals_string.go -type=Dir,Conn,Fmt,Lvl signals_types.go

// Dir represents the signal direction.
type Dir int

const (
	unknownDir Dir = iota
	Input
	Output
)

// Conn represents the physical connection.
type Conn int

const (
	unknownConn Conn = iota
	XLR
	Jack
)

// Fmt represents the signal format.
type Fmt int

const (
	unknownFmt Fmt = iota
	Analog
	AES
)

// Lvl represents the signal level.
type Lvl int

const (
	unknownLvl Lvl = iota
	Line
	Mic
)
