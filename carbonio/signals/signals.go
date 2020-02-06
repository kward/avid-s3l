/*
Package signals provides control over the Carbon I/O signals.

The Avid S3L signals are controlled through the kernel interface of the `spi`
device. The current state of the signal can be read by looking at the contents
of the file that represents the device interface. Changing the state can be done
by writing to the same interface.
*/
package signals

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/kward/avid-s3l/carbonio/helpers"
)

func MicInputs(numInputs uint) (Signals, error) {
	if numInputs == 0 {
		return nil, fmt.Errorf("invalid number of inputs %d", numInputs)
	}

	// Mic inputs. Counting inputs from 1, i.e. 1-16 (not 0-15).
	ss := Signals{}
	for i := uint(1); i <= numInputs; i++ {
		s, err := New(
			fmt.Sprintf("Mic input #%d", i),
			Number(i),
			MaxNumber(numInputs),
			Direction(Input),
			Connector(XLR),
			Format(Analog),
			Level(Mic),
		)
		if err != nil {
			return nil, fmt.Errorf("error instantiating input %d; %s", i, err)
		}
		ss[i] = s
	}

	return ss, nil
}

type Signals map[uint]*Signal

// Signal describes a Carbon I/O signal.
type Signal struct {
	opts *options
	name string

	gainSPI    string
	padSPI     string
	phantomSPI string
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

	var chSPI, chPre string

	s := &Signal{
		opts: o,
		name: name,
	}
	switch o.dir {
	case Input:
		chSPI = channelSPI(o.num, o.maxNum)
		chPre = channelPrefix(o.num, o.maxNum)
	default:
		return nil, fmt.Errorf("unsupported signal direction %s", o.dir)
	}

	s.gainSPI = fmt.Sprintf("%s/%s_preamp_gain", chSPI, chPre)
	s.padSPI = fmt.Sprintf("%s/%s_pad_en", chSPI, chPre)
	s.phantomSPI = phantomSPI(o.num, o.maxNum)

	return s, nil
}

func (s *Signal) Connector() Conn { return s.opts.conn }
func (s *Signal) Format() Fmt     { return s.opts.fmt }
func (s *Signal) Level() Lvl      { return s.opts.lvl }

const gainOffset = 9 // Offset value between spi value and real dB gain.

// Gain returns the current gain level in dB.
//
// The spi gain value is between 1-51, which represents a gain of 10-60 dB.
func (s *Signal) Gain() (uint, error) {
	u, err := readFileUint(s, s.gainSPI)
	if err != nil {
		return 0, err
	}
	if u < 1 || u > 51 {
		return 0, fmt.Errorf("unsupported spi gain value %d", u)
	}
	return u + gainOffset, nil
}

// SetGain for the given signal.
func (s *Signal) SetGain(gain uint) error {
	if gain < 10 || gain > 60 {
		return fmt.Errorf("unsupported gain value %d", gain)
	}
	if err := helpers.WriteFile(s.gainSPI, fmt.Sprintf("%d", gain-gainOffset)); err != nil {
		return fmt.Errorf("error writing gain; %s", err)
	}
	return nil
}

// Pad returns the current state of the -20 dB pad.
func (s *Signal) Pad() (bool, error) {
	v, err := helpers.ReadByte(s.padSPI)
	if err != nil {
		return false, fmt.Errorf("error reading pad; %s", err)
	}
	switch v {
	case '0':
		return false, nil
	case '1':
		return true, nil
	default:
		return false, fmt.Errorf("unsupported spi pad value %d", v)
	}
}

// SetPad for the given signal.
// The pad is controlled with the "spi1.X/chX_pad_en" interface.
func (s *Signal) SetPad(pad bool) error {
	v := '0'
	if pad {
		v = '1'
	}
	return helpers.WriteByte(s.padSPI, byte(v))
}

// Phantom returns the current state of the -48 V phantom.
//
// Phantom states are stored as 4 bit values of a byte, with the lowest signal
// number in the highest bit. The byte itself is stored as a string.
//
// 1 = 8 (0b00001000)
// 2 = 4 (0b00000100)
// 3 = 2 (0b00000010)
// 4 = 1 (0b00000001)
func (s *Signal) Phantom() (bool, error) {
	u, err := readFileUint(s, s.phantomSPI)
	if err != nil {
		return false, err
	}
	if u > 0b00001111 { // The max value when all four agc phantoms are enabled.
		return false, fmt.Errorf("unsupported spi phantom value %d (%08b)", u, u)
	}

	v := uint(1 << (3 - ((s.opts.num - 1) % 4)))
	return u&v > 0, nil
}

// SetPhantom for the given signal.
//
// The phantom is controlled with the "spi4.0/adcX_phantom_en" interface.
func (s *Signal) SetPhantom(phantom bool) error {
	u, err := readFileUint(s, s.phantomSPI)
	if err != nil {
		return err
	}
	if u > 0b00001111 { // The max value when all four agc phantoms are enabled.
		return fmt.Errorf("unsupported spi phantom value %d (%08b)", u, u)
	}

	v := uint(1 << (3 - ((s.opts.num - 1) % 4)))
	if phantom {
		u = u | v
	} else {
		u = u & ^v
	}

	if err := helpers.WriteFile(s.phantomSPI, fmt.Sprintf("%d", u)); err != nil {
		return fmt.Errorf("error writing phantom; %s", err)
	}
	return nil
}

func readFileUint(s *Signal, filename string) (uint, error) {
	data, err := helpers.ReadFileFn(filename)
	if err != nil {
		return 0, fmt.Errorf("error reading phantom; %s", err)
	}
	// Convert data (slice of bytes) to a uint.
	str := strings.Split(fmt.Sprintf("%s", data), "\n")[0]
	u64, err := strconv.ParseUint(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("error converting %q to an int; %s", data, err)
	}
	return uint(u64), nil
}

// channeSPI maps the input number to the appropriate SPI device path.
//
// Input signals are controlled with individual files using the SPI device
// interface. The inputs are spread across devices (e.g., "spi1.1" for input
// signal 1, or "spi1.2" for input signal 16).
//
// See also `channelPrefix()`.
func channelSPI(num, maxNum uint) string {
	const spi = "/sys/bus/spi/devices/spi1."

	if num < 1 || num > maxNum {
		return "unknown"
	}

	switch num {
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

// channelPrefix maps the input number to the appropriate channel prefix.
//
// Input signals are controlled with individual files using the SPI device
// interface. The inputs are spread across devices and channels are prefixed
// with a chX value (e.g., "ch0" for input #1, or "ch3" for input #16).
//
// See also `channelSPI()`.
func channelPrefix(num, maxNum uint) string {
	if num < 1 || num > maxNum {
		return "unknown"
	}
	return fmt.Sprintf("ch%d", (num-1)%4)
}

// phantomSPI maps the input number to the appropriate SPI device path.
//
// Input signals are controlled with individual files using the SPI device
// interface. The inputs are spread across devices and phantoms are prefixed
// with a adcX value (e.g., "adc1" for input #1, or "adc2" for input #16).
//
// See also `channelPrefix()`.
func phantomSPI(num, maxNum uint) string {
	const spi = "/sys/bus/spi/devices/spi4.0/"

	if num < 1 || num > maxNum {
		return "unknown"
	}

	switch num {
	case 1, 2, 3, 4:
		return spi + "adc1_phantom_en"
	case 5, 6, 7, 8:
		return spi + "adc0_phantom_en"
	case 9, 10, 11, 12:
		return spi + "adc3_phantom_en"
	case 13, 14, 15, 16:
		return spi + "adc2_phantom_en"
	default:
		return "unknown"
	}
}
