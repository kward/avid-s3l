/*
Package devices enables control of specific Avid S3L devices.
*/
package devices

import (
	"github.com/kward/avid-s3l/carbonio/leds"
	"github.com/kward/avid-s3l/carbonio/signals"
)

type Device interface {
	// LED returns defined LEDs.
	LEDs() *leds.LEDs
	// NumMicInputs returns the number of microphone inputs for the device.
	NumMicInputs() uint
	// MicInput returns the signal struct for the request input number.
	MicInput(input uint) (*signals.Signal, error)
}

const SPIDevicesDir = "/sys/bus/spi/devices"
