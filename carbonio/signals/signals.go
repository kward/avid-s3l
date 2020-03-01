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
)

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

	gain, err := NewGain(o.num, o.spiDelayRead, o.spiBaseDir)
	if err != nil {
		return nil, err
	}
	pad, err := NewPad(o.num, o.spiDelayRead, o.spiBaseDir)
	if err != nil {
		return nil, err
	}
	phantom, err := NewPhantom(o.num, o.spiDelayRead, o.spiBaseDir)
	if err != nil {
		return nil, err
	}

	s := &Signal{
		opts:    o,
		name:    name,
		gain:    gain,
		pad:     pad,
		phantom: phantom,
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
