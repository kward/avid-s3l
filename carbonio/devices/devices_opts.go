package devices

import (
	"fmt"
	"net"
)

type options struct {
	// Device options.
	mac  net.HardwareAddr
	ip   net.IP
	host string

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

func Host(v string) func(*options) error {
	return func(o *options) error { return o.setHost(v) }
}
func (o *options) setHost(v string) error {
	o.host = v
	return nil
}

func IP(v net.IP) func(*options) error {
	return func(o *options) error { return o.setIP(v) }
}
func (o *options) setIP(v net.IP) error {
	o.ip = v
	return nil
}

func MAC(v net.HardwareAddr) func(*options) error {
	return func(o *options) error { return o.setMAC(v) }
}
func (o *options) setMAC(v net.HardwareAddr) error {
	o.mac = v
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
