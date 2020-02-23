package leds

type options struct {
	// SPI options.
	spiDelayRead bool // Delay SPI Read() until first direct call.
	// Global flags.
	spiBaseDir string
	verbose    bool
}

func (o *options) validate() error {
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
