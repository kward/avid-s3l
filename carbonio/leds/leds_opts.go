package leds

import "fmt"

type options struct {
	// Global flags.
	spiBaseDir string
	verbose    bool
}

func (o *options) validate() error {
	if o.spiBaseDir == "" {
		return fmt.Errorf("SPIBaseDir option missing")
	}
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
