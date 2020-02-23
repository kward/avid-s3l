package signals

import "fmt"

type options struct {
	// Device options.
	maxNum int // Maximum number of signals of this type.
	// Signal options.
	num  int // Signal number (1-based, i.e. 1 is the first signal number).
	conn ConnectorEnum
	dir  DirectionEnum
	fmt  FormatEnum
	lvl  LevelEnum
	// SPI options.
	spiDelayRead bool // Delay SPI Read() until first direct call.
	// Global flags.
	spiBaseDir string
	verbose    bool
}

func (o *options) validate() error {
	if o.maxNum == 0 {
		return fmt.Errorf("MaxNumber option missing")
	}
	if o.num <= 0 || o.num > o.maxNum {
		return fmt.Errorf("Number option out of range [0:%d]", o.maxNum)
	}
	return nil
}

// MaxNumber returns the maximum number of supported signals.
func MaxNumber(v int) func(*options) error {
	return func(o *options) error { return o.setMaxNumber(v) }
}
func (o *options) setMaxNumber(v int) error {
	o.maxNum = v
	return nil
}

// Number returns the signal number (1-based).
func Number(v int) func(*options) error {
	return func(o *options) error { return o.setNumber(v) }
}
func (o *options) setNumber(v int) error {
	o.num = v
	return nil
}

// Connector returns the type of connector for the signal.
func Connector(v ConnectorEnum) func(*options) error {
	return func(o *options) error { return o.setConnector(v) }
}
func (o *options) setConnector(v ConnectorEnum) error {
	o.conn = v
	return nil
}

// Direction returns the signal direction.
func Direction(v DirectionEnum) func(*options) error {
	return func(o *options) error { return o.setDirection(v) }
}
func (o *options) setDirection(v DirectionEnum) error {
	o.dir = v
	return nil
}

// Format returns the format of the signal.
func Format(v FormatEnum) func(*options) error {
	return func(o *options) error { return o.setFormat(v) }
}
func (o *options) setFormat(v FormatEnum) error {
	o.fmt = v
	return nil
}

// Level returns the level of the signal.
func Level(v LevelEnum) func(*options) error {
	return func(o *options) error { return o.setLevel(v) }
}
func (o *options) setLevel(v LevelEnum) error {
	o.lvl = v
	return nil
}

// SPIDelayRead returns whether SPI Read() should be delayed until first call.
func SPIDelayRead(v bool) func(*options) error {
	return func(o *options) error { return o.setSPIDelayRead(v) }
}
func (o *options) setSPIDelayRead(v bool) error {
	o.spiDelayRead = v
	return nil
}

// SPIBaseDir returns the path to the SPI devices directory.
func SPIBaseDir(v string) func(*options) error {
	return func(o *options) error { return o.setSPIBaseDir(v) }
}
func (o *options) setSPIBaseDir(v string) error {
	o.spiBaseDir = v
	return nil
}

// Verbose returns the verbose output setting.
func Verbose(v bool) func(*options) error {
	return func(o *options) error { return o.setVerbose(v) }
}
func (o *options) setVerbose(v bool) error {
	o.verbose = v
	return nil
}
