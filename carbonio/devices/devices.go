/*
Package devices enables control of specific Avid S3L devices.
*/
package devices

import (
	"net"

	"github.com/kward/avid-s3l/carbonio/signals"
)

type Device interface {
	// NumMicInputs returns the number of microphone inputs for the device.
	NumMicInputs() uint
	// MicInput returns the signal struct for the request input number.
	MicInput(input uint) (*signals.Signal, error)
}

const SPIDevicesDir = "/sys/bus/spi/devices"

type options struct {
	mac  net.HardwareAddr
	ip   net.IP
	host string

	spiBaseDir string
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
