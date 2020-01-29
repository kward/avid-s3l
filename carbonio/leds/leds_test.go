package leds

import (
	"fmt"
	"os"
	"testing"

	"github.com/kward/golib/operators"
)

type testCase struct {
	desc  string
	state LEDState
	data  []byte
	ok    bool
}

var (
	testCases = []testCase{
		// Common states.
		{"off", Off, []byte{'0', '\n'}, true},
		// Unknown states.
		{"unknown", Unknown, []byte{123, '\n'}, false},
		// Data errors.
		{"zero length data", Unknown, []byte{}, false},
		{"too much data", Unknown, []byte{1, 2, 3, 4}, false},
		{"wrong termination", Unknown, []byte{'2', 34}, false},
		// ReadFile errors.
		{"readfile error", Unknown, []byte{}, false},
	}
)

func TestPowerLED(t *testing.T) {
	led := new(powerLED)

	for _, tc := range append(testCases, []testCase{
		{"on", On, []byte{'2', '\n'}, true},
		{"alert", Alert, []byte{'1', '\n'}, true},
	}...) {
		t.Run(fmt.Sprintf("State() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.data, tc.ok)
			got, err := led.State()
			if err != nil && tc.ok == true {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && tc.ok == false {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if want := tc.state; got != want {
				t.Errorf("= %s, want %s", got, want)
			}
		})

		t.Run(fmt.Sprintf("SetState() %s", tc.desc), func(t *testing.T) {
			prepareWriteFile(tc.ok)
			err := led.SetState(tc.state)
			if err != nil && tc.ok == true {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && tc.ok == false {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := wfData, tc.data; !operators.EqualSlicesOfByte(got, want) {
				t.Errorf("expected %v to be written, not %v", want, got)
			}
		})
	}
}

func TestStatusLED(t *testing.T) {
	led := new(statusLED)
	for _, tc := range append(testCases, []testCase{
		{"on", On, []byte{'2', '\n'}, true},
		{"alert", Alert, []byte{'1', '\n'}, true},
	}...) {
		t.Run(fmt.Sprintf("State() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.data, tc.ok)
			got, err := led.State()
			if err != nil && tc.ok == true {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && tc.ok == false {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if want := tc.state; got != want {
				t.Errorf("State() = %s, want %s", got, want)
			}
		})

		t.Run(fmt.Sprintf("SetState() %s", tc.desc), func(t *testing.T) {
			prepareWriteFile(tc.ok)
			err := led.SetState(tc.state)
			if err != nil && tc.ok == true {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && tc.ok == false {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := wfData, tc.data; !operators.EqualSlicesOfByte(got, want) {
				t.Errorf("expected %v to be written, not %v", want, got)
			}
		})
	}
}

func TestMuteLED(t *testing.T) {
	led := new(muteLED)
	for _, tc := range append(testCases, []testCase{
		{"on", On, []byte{'1', '\n'}, true},
	}...) {
		t.Run(fmt.Sprintf("State() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.data, tc.ok)
			got, err := led.State()
			if err != nil && tc.ok == true {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && tc.ok == false {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if want := tc.state; got != want {
				t.Errorf("State() = %s, want %s", got, want)
			}
		})

		t.Run(fmt.Sprintf("SetState() %s", tc.desc), func(t *testing.T) {
			prepareWriteFile(tc.ok)
			err := led.SetState(tc.state)
			if err != nil && tc.ok == true {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && tc.ok == false {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := wfData, tc.data; !operators.EqualSlicesOfByte(got, want) {
				t.Errorf("expected %v to be written, not %v", want, got)
			}
		})
	}
}

func TestMain(m *testing.M) {
	// Override function pointers in leds for testing.
	readFileFn = readFile
	writeFileFn = writeFile

	os.Exit(m.Run())
}

var (
	rfData, wfData []byte
	rfErr, wfErr   error
)

func prepareReadFile(data []byte, ok bool) {
	rfData = data
	if !ok {
		rfErr = fmt.Errorf("ReadFile error for testing")
	}
	rfErr = nil
}

// readFile matches the signature of io.ReadFile.
func readFile(filename string) ([]byte, error) {
	return rfData, rfErr
}

func prepareWriteFile(ok bool) {
	if !ok {
		wfErr = fmt.Errorf("WriteFile error for testing")
	}
	wfErr = nil
}

// writeFile matches the signature of io.WriteFile.
func writeFile(filename string, data []byte, mode os.FileMode) error {
	if wfErr != nil {
		wfData = []byte{}
		return wfErr
	}
	wfData = data
	return nil
}
