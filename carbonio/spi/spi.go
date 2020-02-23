/*
Package spi reads-and-writes data directly to-and-from the SPI files. There is
no processing of the data beyond converting it between []byte and int types. Any
higher level understanding of the data (e.g., pad bitwise values) must be
implemented downstream.
*/
package spi

import (
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/kward/avid-s3l/carbonio/helpers"
)

const DevicesDir = "/sys/bus/spi/devices"

type SPI struct {
	opts *options

	enum  Enum
	path  string // Full path to the SPI file (e.g. `/some/path/spi1.0/ch0_pad_en`).
	value int    // Current value of the SPI file.
	raw   []byte // Most recent raw value read-or-written.
}

func New(enum Enum, num int, opts ...func(*options) error) (*SPI, error) {
	o := &options{}
	o.setBaseDir(DevicesDir)
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}
	if err := o.validate(); err != nil {
		return nil, err
	}

	spi := &SPI{
		opts: o,
		enum: enum,
	}
	var p string
	switch enum {
	// LEDs.
	case Power:
		p = "spi4.0/status_led_1_en"
	case Status:
		p = "spi4.0/status_led_0_en"
	case Mute:
		p = "spi4.0/mute_led_en"
	// Inputs.
	case Gain:
		p = path.Join(channelDir(num), fmt.Sprintf("ch%d_preamp_gain", channelNum(num)))
	case Pad:
		p = path.Join(channelDir(num), fmt.Sprintf("ch%d_pad_en", channelNum(num)))
	case Phantom:
		p = path.Join("spi4.0", phantomFile(num))
		// Testing (for testing only).
	case Blinky:
		p = "spiX.Y/blinky_en"
	}
	spi.path = path.Join(o.baseDir, p)

	if !o.delayRead {
		if _, err := spi.Read(); err != nil {
			return nil, err
		}
	}
	return spi, nil
}

// Read the current value from the SPI interface, storing a copy in `value`.
func (s *SPI) Read() (int, error) {
	data, err := helpers.ReadFileFn()(s.path)
	if err != nil {
		return 0, fmt.Errorf("failed to read %s from %s; %s", s.enum, s.path, err)
	}
	if len(data) == 0 {
		return 0, fmt.Errorf("empty data read from %s", s.path)
	}
	s.raw = data
	v, err := strconv.Atoi(strings.TrimRight(string(data), "\n"))
	if err != nil {
		return 0, fmt.Errorf("conversion failure; %s", err)
	}
	s.value = v
	return v, nil
}

var fileMode = os.FileMode(0644)

// Write data to the SPI interface.
func (s *SPI) Write(v int) error {
	str := strconv.Itoa(v) + "\n"
	if err := helpers.WriteFileFn()(s.path, []byte(str), fileMode); err != nil {
		return fmt.Errorf("failed to write %s value of %d to %s; %s", s.enum, v, s.path, err)
	}

	// Do a read-after-write to verify the data, forcing a data re-population.
	data, err := s.Read()
	if err != nil {
		return fmt.Errorf("read-after-write error: %s", err)
	}
	if v != data {
		return fmt.Errorf("read-after-write data mismatch: got = %d, want = %d", data, v)
	}
	return nil
}

// Implementation is the interface that describes a minimal implementation.
type Implementation interface {
	// Initialize the SPI device.
	Initialize() error
	// Name returns the SPI name.
	Name() string
	// Path in the file system to the SPI device.
	Path() string
	// Raw returns the most recent raw value read from the SPI interface.
	Raw() []byte
}

// InitializeFn describes the Initialize function.
type InitializeFn func() error

// Initialize implements Implementation.
func (s *SPI) Initialize() error { return s.Write(0) }

// Name implements Implementation.
func (s *SPI) Name() string { return s.enum.String() }

// PathFn describes the Path funcion.
type PathFn func() string

// Path implements Implementation.
func (s *SPI) Path() string { return string(s.path) }

// Raw implements Implementation.
func (s *SPI) Raw() []byte { return s.raw }

// Value returns the most recent value read from the SPI interface.
func (s *SPI) Value() int { return s.value }

// channelDir maps the input number to the appropriate SPI device directory.
//
// Input signals are controlled with individual files using the SPI device
// interface. The inputs are spread across devices (e.g., `spi1.1` for input
// signal 1, or `spi1.2` for input signal 16).
//
// See also `channelNum()`.
func channelDir(num int) string {
	switch num {
	case 1, 2, 3, 4:
		return "spi1.1"
	case 5, 6, 7, 8:
		return "spi1.0"
	case 9, 10, 11, 12:
		return "spi1.3"
	case 13, 14, 15, 16:
		return "spi1.2"
	default:
		return "unknown"
	}
}

// channelNum maps the input number to the appropriate channel number.
//
// Input signals are controlled with individual files using the SPI device
// interface. The inputs are spread across devices and channels are prefixed
// with a chX value (e.g., `ch0` for input #1, or `ch3` for input #16).
//
// See also `channelDir()`.
func channelNum(num int) int {
	return (num - 1) % 4
}

// phantomFile maps the input number to the appropriate SPI device file.
//
// Input signals are controlled with individual files using the SPI device
// interface. The inputs are spread across devices and phantoms are prefixed
// with a adcX value (e.g., `adc1` for input #1, or `adc2` for input #16).
func phantomFile(num int) string {
	var spi string
	switch num {
	case 1, 2, 3, 4:
		spi = "adc1_phantom_en"
	case 5, 6, 7, 8:
		spi = "adc0_phantom_en"
	case 9, 10, 11, 12:
		spi = "adc3_phantom_en"
	case 13, 14, 15, 16:
		spi = "adc2_phantom_en"
	default:
		return "unknown"
	}
	return spi
}
