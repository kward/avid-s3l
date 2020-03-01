package signals

import "fmt"

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
