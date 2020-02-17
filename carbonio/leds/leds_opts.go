package leds

import "fmt"

type options struct {
	// Global flags.
	spiBaseDir string
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
