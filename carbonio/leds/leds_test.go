package leds

import (
	"fmt"
	"os"
	"testing"

	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/golib/operators"
)

func TestMain(m *testing.M) {
	// Override function pointers for testing.
	helpers.SetReadFileFn(helpers.MockReadFile)
	helpers.SetWriteFileFn(helpers.MockWriteFile)

	os.Exit(m.Run())
}

func TestLEDs(t *testing.T) {
	_, err := New(SPIBaseDir("/spi/base"))
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}
}

func TestLED(t *testing.T) {
	led, err := NewLED("Blinky", "/path/to/blinky",
		byState{Off: '0', Alert: '1', On: '2', testState: 255},
		SPIBaseDir("/spi/base"),
	)
	if err != nil {
		t.Fatalf("unexpected error; %s", err)
	}

	for _, tc := range []struct {
		desc  string
		state State
		data  []byte
		ok    bool
		rfErr error
		wfErr error
	}{
		// Supported states.
		{"off", Off, []byte{'0', '\n'}, true, nil, nil},
		{"on", On, []byte{'2', '\n'}, true, nil, nil},
		{"alert", Alert, []byte{'1', '\n'}, true, nil, nil},

		// Unknown states.
		{desc: "unknown", data: []byte{123, '\n'}},
		// Data errors.
		{desc: "zero length data"},
		{desc: "too much data", data: []byte{1, 2, 3, 4}},
		{desc: "wrong termination", data: []byte{'2', 34}},
		// ReadFile errors.
		{desc: "readfile error", rfErr: fmt.Errorf("ReadFile error")},
		// Write errors.
		{desc: "writefile error", wfErr: fmt.Errorf("WriteFile error")},
	} {
		t.Run(fmt.Sprintf("State() %s", tc.desc), func(t *testing.T) {
			helpers.PrepareReadFile(tc.data, tc.rfErr)
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

		t.Run(fmt.Sprintf("SetState() %s", tc.desc), func(t *testing.T) {
			helpers.PrepareWriteFile(tc.wfErr)
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
			if got, want := helpers.MockWriteData(), tc.data; !operators.EqualSlicesOfByte(got, want) {
				t.Errorf("expected %v to be written, not %v", want, got)
			}
		})

		t.Run(fmt.Sprintf("Path() %s", tc.desc), func(t *testing.T) {
			if !tc.ok {
				return
			}
			if got, want := led.Path(), "/path/to/blinky"; got != want {
				t.Errorf("= %s, want %s", got, want)
			}
		})
	}
}
