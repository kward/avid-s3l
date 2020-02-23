package leds

import (
	"fmt"
	"os"
	"testing"

	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/avid-s3l/carbonio/spi"
)

func TestMain(m *testing.M) {
	// Override function pointers for testing.
	helpers.SetReadFileFn(helpers.MockReadFile)
	helpers.SetWriteFileFn(helpers.MockWriteFile)

	os.Exit(m.Run())
}

func TestLEDs(t *testing.T) {
	_, err := New(SPIDelayRead(true))
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}
}

func newLED() (*LED, error) {
	return NewLED(spi.Blinky,
		byState{Off: 0, Alert: 1, On: 2, testState: 255},
		SPIDelayRead(true),
	)
}

func TestState(t *testing.T) {
	led, err := newLED()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}

	for _, tc := range []struct {
		desc     string
		ok       bool
		rfErr    error
		spiValue int // Current SPI value.

		state State
	}{
		// Supported states.
		{"off", true, nil, 0, Off},
		{"alert", true, nil, 1, Alert},
		{"on", true, nil, 2, On},

		// Error states.
		{desc: "unsupported spi value", spiValue: 123},
		{desc: "readfile error", rfErr: fmt.Errorf("mock ReadFile error")},
	} {
		t.Run(fmt.Sprintf("State() %s", tc.desc), func(t *testing.T) {
			helpers.ResetMockReadWrite()
			helpers.PrepareMockReadFile([]byte{}, tc.rfErr)
			led.spi.Write(tc.spiValue)

			got, err := led.State()
			if err != nil && tc.ok {
				t.Fatalf("unexpected error; %s", err)
			}
			if err == nil && !tc.ok {
				t.Fatal("expected an error")
			}
			if !tc.ok {
				return
			}
			if want := tc.state; got != want {
				t.Errorf("= %s, want %s", got, want)
			}
		})
	}
}

func TestSetState(t *testing.T) {
	led, err := newLED()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}

	for _, tc := range []struct {
		desc     string
		ok       bool
		wfErr    error
		spiValue int // Current SPI value.

		state State
		value int
	}{
		// Supported states.
		{"off to off", true, nil, 0, Off, 0},
		{"off to alert", true, nil, 0, Alert, 1},
		{"off to on ", true, nil, 0, On, 2},
		{"on to on ", true, nil, 2, On, 2},
		{"on to alert", true, nil, 2, Alert, 1},
		{"on to off", true, nil, 2, Off, 0},

		// Unknown states.
		{desc: "unsupported spi value", spiValue: 123},
		{desc: "writefile error", wfErr: fmt.Errorf("mock WriteFile error")},
	} {
		t.Run(fmt.Sprintf("SetState() %s", tc.desc), func(t *testing.T) {
			helpers.ResetMockReadWrite()
			helpers.PrepareMockWriteFile(tc.wfErr)
			led.spi.Write(tc.spiValue)

			err := led.SetState(tc.state)
			if err != nil && tc.ok {
				t.Fatalf("unexpected error; %s", err)
			}
			if err == nil && !tc.ok {
				t.Fatal("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := led.spi.Value(), tc.value; got != want {
				t.Errorf("SPI Value() = %d, want %d", want, got)
			}
		})
	}
}

func TestName(t *testing.T) {
	led, err := newLED()
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}

	t.Run("Name()", func(t *testing.T) {
		if got, want := led.Name(), spi.Blinky.String(); got != want {
			t.Errorf("= %q, want %q", got, want)
		}
	})
}
