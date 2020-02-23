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
	"os"

	"github.com/kward/avid-s3l/carbonio/spi"
)

func MicInputs(numInputs int, opts ...func(*options) error) (Signals, error) {
	if numInputs == 0 {
		return nil, fmt.Errorf("invalid number of inputs %d", numInputs)
	}

	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	// Mic inputs. Counting inputs from 1, i.e. 1-16 (not 0-15).
	ss := Signals{}
	for i := 1; i <= numInputs; i++ {
		s, err := New(
			fmt.Sprintf("Mic input #%d", i),
			MaxNumber(numInputs),
			Number(i),
			Direction(Input),
			Connector(XLR),
			Format(Analog),
			Level(Mic),
			SPIDelayRead(o.spiDelayRead),
			SPIBaseDir(o.spiBaseDir),
			Verbose(o.verbose),
		)
		if err != nil {
			return nil, fmt.Errorf("error instantiating input %d; %s", i, err)
		}
		ss[i] = s
	}

	return ss, nil
}

type Signals map[int]*Signal

// Signal describes a Carbon I/O signal.
type Signal struct {
	opts *options

	name    string
	gain    *Gain
	pad     *Pad
	phantom *Phantom
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

	gain, err := spi.New(spi.Gain, o.num,
		spi.DelayRead(o.spiDelayRead),
		spi.BaseDir(o.spiBaseDir),
	)
	if err != nil {
		return nil, fmt.Errorf("failure instantiating Gain SPI; %s", err)
	}
	pad, err := spi.New(spi.Pad, o.num,
		spi.DelayRead(o.spiDelayRead),
		spi.BaseDir(o.spiBaseDir),
	)
	if err != nil {
		return nil, fmt.Errorf("failure instantiating Pad SPI; %s", err)
	}
	phantom, err := spi.New(spi.Phantom, o.num,
		spi.DelayRead(o.spiDelayRead),
		spi.BaseDir(o.spiBaseDir),
	)
	if err != nil {
		return nil, fmt.Errorf("failure instantiating Phantom SPI; %s", err)
	}

	s := &Signal{
		opts:    o,
		name:    name,
		gain:    &Gain{gain},
		pad:     &Pad{pad},
		phantom: &Phantom{phantom, o.num},
	}

	if o.verbose {
		fmt.Fprintf(os.Stderr, "%#v\n", s)
	}
	return s, nil
}

func (s *Signal) Connector() ConnectorEnum { return s.opts.conn }
func (s *Signal) Direction() DirectionEnum { return s.opts.dir }
func (s *Signal) Format() FormatEnum       { return s.opts.fmt }
func (s *Signal) Level() LevelEnum         { return s.opts.lvl }

func (s *Signal) Gain() *Gain       { return s.gain }
func (s *Signal) Pad() *Pad         { return s.pad }
func (s *Signal) Phantom() *Phantom { return s.phantom }

// Gain provides access to the gain SPI.
type Gain struct {
	spi *spi.SPI
}

// Ensure spi interfaces are implemented.
var _ spi.Implementation = new(Gain)

const (
	gainMin    = 10
	gainMax    = 60
	gainOffset = 9 // Offset between SPI value and real dB gain.
)

// Value returns the gain level in dB.
//
// The SPI gain value is between 1-51, which represents a gain of 10-60 dB.
func (g *Gain) Value() (uint, error) {
	v, err := g.spi.Read()
	if err != nil {
		return 0, fmt.Errorf("error reading gain; %s", err)
	}
	if v < (gainMin-gainOffset) || v > (gainMax-gainOffset) {
		return 0, fmt.Errorf("unsupported spi gain value %d", v)
	}
	return uint(v) + gainOffset, nil
}

// SetValue of gain in dB.
func (g *Gain) SetValue(gain uint) error {
	if gain < gainMin || gain > gainMax {
		return fmt.Errorf("unsupported gain value %d", gain)
	}
	if err := g.spi.Write(int(gain - gainOffset)); err != nil {
		return fmt.Errorf("error writing gain; %s", err)
	}
	return nil
}

// Initialize implements spi.Implementation.
func (g *Gain) Initialize() error { return g.SetValue(gainMin) }

// Name implements spi.Implementation.
func (g *Gain) Name() string { return g.spi.Name() }

// Path implements spi.Implementation.
func (g *Gain) Path() string { return g.spi.Path() }

// Raw implements spi.Implementation.
func (g *Gain) Raw() []byte { return g.spi.Raw() }

// Pad provides access to the pad SPI.
type Pad struct {
	spi *spi.SPI
}

// Ensure spi interfaces are implemented.
var _ spi.Implementation = new(Pad)

const (
	PadEnabled  = true
	PadDisabled = false
)

// Enable the pad.
func (p *Pad) Enable() error {
	return p.setState(PadEnabled)
}

// Disable the pad.
func (p *Pad) Disable() error {
	return p.setState(PadDisabled)
}

func (p *Pad) setState(state bool) error {
	v := 0
	if state {
		v = 1
	}
	if err := p.spi.Write(v); err != nil {
		return fmt.Errorf("error writing pad; %s", err)
	}
	return nil
}

// IsEnabled returns whether the -20 dB pad is enabled.
func (p *Pad) IsEnabled() (bool, error) {
	v, err := p.spi.Read()
	if err != nil {
		return false, fmt.Errorf("error reading pad; %s", err)
	}
	switch v {
	case 0:
		return PadDisabled, nil
	case 1:
		return PadEnabled, nil
	default:
		return false, fmt.Errorf("unsupported spi pad value %d", v)
	}
}

// Initialize implements spi.Implementation.
func (p *Pad) Initialize() error { return p.Disable() }

// Name implements spi.Implementation.
func (p *Pad) Name() string { return p.spi.Name() }

// Path implements spi.Implementation.
func (p *Pad) Path() string { return p.spi.Path() }

// Raw implements spi.Implementation.
func (p *Pad) Raw() []byte { return p.spi.Raw() }

// Phantom provides access to the phantom SPI.
type Phantom struct {
	spi *spi.SPI
	num int
}

// Ensure spi interfaces are implemented.
var _ spi.Implementation = new(Phantom)

const (
	PhantomEnabled  = true
	PhantomDisabled = false
)

// Enable the phantom.
func (p *Phantom) Enable() error {
	return p.setState(PhantomEnabled)
}

// Disable the phantom.
func (p *Phantom) Disable() error {
	return p.setState(PhantomDisabled)
}

func (p *Phantom) setState(state bool) error {
	u := uint(p.spi.Value())
	v := uint(1 << (3 - ((p.num - 1) % 4)))
	if state == PhantomEnabled {
		u = u | v
	} else {
		u = u & ^v
	}
	if err := p.spi.Write(int(u)); err != nil {
		return fmt.Errorf("error writing phantom; %s", err)
	}
	return nil
}

// IsEnabled returns whether the -48 V phantom is enabled.
//
// Phantom states are stored as 4 bit values of a byte, with the lowest signal
// number in the highest bit. The byte itself is stored as a string.
//
// 1 = 8 (0b00001000)
// 2 = 4 (0b00000100)
// 3 = 2 (0b00000010)
// 4 = 1 (0b00000001)
func (p *Phantom) IsEnabled() (bool, error) {
	v, err := p.spi.Read()
	if err != nil {
		return false, fmt.Errorf("error reading phantom; %s", err)
	}
	if v > 0b00001111 { // The max value when all four agc phantoms are enabled.
		return false, fmt.Errorf("unsupported spi phantom value %d (%08b)", v, v)
	}
	u := uint(v)
	w := uint(1 << (3 - ((p.num - 1) % 4)))
	return u&w > 0, nil
}

// Initialize implements spi.Implementation.
func (p *Phantom) Initialize() error { return p.Disable() }

// Name implements spi.Implementation.
func (p *Phantom) Name() string { return p.spi.Name() }

// Path implements spi.Implementation.
func (p *Phantom) Path() string { return p.spi.Path() }

// Raw implements spi.Implementation.
func (p *Phantom) Raw() []byte { return p.spi.Raw() }
