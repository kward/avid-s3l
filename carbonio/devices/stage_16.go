package devices

import (
	"fmt"

	"github.com/kward/avid-s3l/carbonio/signals"
)

const (
	numMicInputs   = 16
	numLineOutputs = 8
	numAESOutputs  = 2
)

type Stage16 struct {
	opts *options

	micInputs map[int]*signals.Signal
}

// Verify that the interface is implemented properly.
var _ Device = new(Stage16)

// NewStage16 returns a populated Stage16 struct.
func NewStage16(opts ...func(*options) error) (*Stage16, error) {
	o := &options{}
	for _, opt := range opts {
		if err := opt(o); err != nil {
			return nil, err
		}
	}

	d := &Stage16{
		opts:      o,
		micInputs: map[int]*signals.Signal{},
	}

	// Counting inputs from 1, i.e. 1-16 (not 0-15).
	for i := 1; i <= numMicInputs; i++ {
		input, err := signals.New(
			fmt.Sprintf("Mic input #%d", i),
			signals.Number(i),
			signals.MaxNumber(numMicInputs),
			signals.Direction(signals.Input),
			signals.Connector(signals.XLR),
			signals.Format(signals.Analog),
			signals.Level(signals.Mic),
		)
		if err != nil {
			return nil, fmt.Errorf("error instantiating input %d; %s", i, err)
		}
		d.micInputs[i] = input
	}

	return d, nil
}
