package signals

import (
	"fmt"
	"testing"

	"github.com/kward/avid-s3l/carbonio/spi"
)

func TestNewInput(t *testing.T) {
	for _, tc := range []struct {
		desc string
		ok   bool

		num    int
		maxNum int
	}{
		// Valid inputs.
		{"signal #1", true, 1, 16},
		{"signal #2", true, 2, 16},
		{"signal #15", true, 15, 16},
		{"signal #16", true, 16, 16},

		// Invalid inputs.
		{desc: "Number out of range", num: 99, maxNum: 16},
		{desc: "MaxNumber not set", num: 16, maxNum: 0},
	} {
		t.Run(fmt.Sprintf("New() %s", tc.desc), func(t *testing.T) {
			s, err := newInput("Beep-ba-beep", tc.num, tc.maxNum)
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatal("expected an error")
			}
			if !tc.ok {
				return
			}
			// Verify that the SPI interfaces were setup.
			if got, want := s.gain.Name(), spi.Gain.String(); got != want {
				t.Errorf("gain Name() = %q, want %q", got, want)
			}
			if got, want := s.pad.Name(), spi.Pad.String(); got != want {
				t.Errorf("pad Name() = %q, want %q", got, want)
			}
			if got, want := s.phantom.Name(), spi.Phantom.String(); got != want {
				t.Errorf("phantom Name() = %q, want %q", got, want)
			}
		})
	}
}

func newInput(name string, num, maxNum int) (*Signal, error) {
	return New(name,
		Number(num),
		MaxNumber(maxNum),
		Direction(Input),
		SPIDelayRead(true), // Prevent initial read from unprepared SPI.
	)
}
