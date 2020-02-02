// Package signals provides control over the Carbon I/O signals.
//
// The Avid S3L signals are controlled through the kernel interface of the spi
// device. The current state of the signal can be read by looking at the
// contents of the file that represents the device interface. Changing the state
// can be done by writing to the same interface.
package signals

import (
	"fmt"

	"github.com/kward/avid-s3l/carbonio/helpers"
)

const (
	NumMicInputs   = 16
	NumLineOutputs = 8
	NumAESOutputs  = 2
)

var (
	MicInputs   []Signal
	LineOutputs []Signal
	AESOutputs  []Signal
)

func init() {

}

// Signal describes a Carbon I/O signal.
type Signal struct {
	opts *options
	name string

	gain    int
	pad     bool
	padFile string
	phantom bool

	iface   string
	channel string
}

// New instantiates a new Signal.
func New(name string, opts ...func(*options) error) (*Signal, error) {
	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	if err := o.validate(); err != nil {
		return nil, err
	}

	signal := &Signal{
		opts: o,
		name: name,
	}
	switch o.dir {
	case Input:
		signal.iface = inputDevicePath(o.num)
		signal.channel = inputChannelPrefix(o.num)
	default:
		return nil, fmt.Errorf("unsupported signal direction %s", o.dir)
	}

	signal.padFile = fmt.Sprintf("%s/%s_pad_en", signal.iface, signal.channel)

	return signal, nil
}

func (s *Signal) Connector() Conn { return s.opts.conn }
func (s *Signal) Format() Fmt     { return s.opts.fmt }
func (s *Signal) Level() Lvl      { return s.opts.lvl }

// Pad returns the state of the -20 dB pad.
func (s *Signal) Pad() (bool, error) {
	v, err := helpers.ReadByte(s.padFile)
	if err != nil {
		return false, fmt.Errorf("error reading pad; %s", err)
	}
	switch v {
	case '0':
		return false, nil
	case '1':
		return true, nil
	default:
		return false, fmt.Errorf("unsupported pad value %d", v)
	}
}

// SetPad for the given signal.
// The pad is controlled with the "spi1.X/chX_pad_en" interface.
func (s *Signal) SetPad(pad bool) error {
	v := '0'
	if pad {
		v = '1'
	}
	return helpers.WriteByte(s.padFile, byte(v))
}

func (s *Signal) Phantom() (bool, error) {
	return false, fmt.Errorf("unimplemented")
}

func (s *Signal) SetPhantom(v bool) error {
	return fmt.Errorf("unimplemented")
}

// inputDevicePath maps the input signal to the appropriate SPI device path.
//
// Input signals are controlled with individual files using the SPI device
// interface. The inputs are spread across devices (e.g., "spi1.1" for input
// signal 1, or "spi1.2" for input signal 16).
//
// See also `inputChannelPrefix()`.
func inputDevicePath(i int) string {
	const spi = "/sys/bus/spi/devices/spi1."
	switch i {
	case 1, 2, 3, 4:
		return spi + "1"
	case 5, 6, 7, 8:
		return spi + "0"
	case 9, 10, 11, 12:
		return spi + "3"
	case 13, 14, 15, 16:
		return spi + "2"
	default:
		return "unknown"
	}
}

// inputChannelPrefix maps the input signal to the appropriate channel prefix.
//
// Input signals are controlled with individual files using the SPI device
// interface. The inputs are spread across devices and are prefixed with a chX
// value (e.g., "ch0" for input signal 1, or "ch3" for input signal 16).
//
// See also `inputDevicePath()`.
func inputChannelPrefix(i int) string {
	if i < 1 || i > NumMicInputs {
		return "unknown"
	}
	return fmt.Sprintf("ch%d", (i-1)%4)
}
