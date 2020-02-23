package spi

type options struct {
	// SPI options.
	delayRead bool // Delay Read() until first call.

	// Global flags.
	baseDir string // spiBaseDir
	verbose bool
}

func (o *options) validate() error {
	return nil
}

// DelayRead returns whether Read() should be delayed until first call.
func DelayRead(v bool) func(*options) error {
	return func(o *options) error { return o.setDelayRead(v) }
}
func (o *options) setDelayRead(v bool) error {
	o.delayRead = v
	return nil
}

// BaseDir returns the path to the SPI devices directory.
func BaseDir(v string) func(*options) error {
	return func(o *options) error { return o.setBaseDir(v) }
}
func (o *options) setBaseDir(v string) error {
	o.baseDir = v
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
