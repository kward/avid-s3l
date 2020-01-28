package leds

import (
	"fmt"
	"os"
	"testing"
)

type testCase struct {
	desc  string
	data  []byte
	ok    bool
	state LEDState
}

var testCases = []testCase{
	// Common states.
	{"off", []byte{48, 10}, true, Off}, // "0"
	// Unknown states.
	{"unknown", []byte{123, 10}, false, Unknown},
	// Data errors.
	{"zero length data", []byte{}, false, Unknown},
	{"too much data", []byte{1, 2, 3, 4}, false, Unknown},
	{"wrong termination", []byte{12, 34}, false, Unknown},
	// ReadFile errors.
	{"readfile error", []byte{}, false, Unknown},
}

func TestPowerLED(t *testing.T) {
	led := new(powerLED)
	for _, tc := range append(testCases, []testCase{
		{"on", []byte{50, 10}, true, On},       // "2"
		{"alert", []byte{49, 10}, true, Alert}, // "1"
	}...) {
		t.Run(fmt.Sprintf("%s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.data, tc.ok)
			got, err := led.State()
			if err != nil && tc.ok == true {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && tc.ok == false {
				t.Fatalf("expected an error")
			}
			if want := tc.state; got != want {
				t.Errorf("State() = %s, want %s", got, want)
			}
		})
	}
}

func TestStatusLED(t *testing.T) {
	led := new(statusLED)
	for _, tc := range append(testCases, []testCase{
		{"on", []byte{50, 10}, true, On},       // "2"
		{"alert", []byte{49, 10}, true, Alert}, // "1"
	}...) {
		t.Run(fmt.Sprintf("%s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.data, tc.ok)
			got, err := led.State()
			if err != nil && tc.ok == true {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && tc.ok == false {
				t.Fatalf("expected an error")
			}
			if want := tc.state; got != want {
				t.Errorf("State() = %s, want %s", got, want)
			}
		})
	}
}

func TestMuteLED(t *testing.T) {
	led := new(muteLED)
	for _, tc := range append(testCases, []testCase{
		{"on", []byte{49, 10}, true, On}, // "1"
	}...) {
		t.Run(fmt.Sprintf("%s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.data, tc.ok)
			got, err := led.State()
			if err != nil && tc.ok == true {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && tc.ok == false {
				t.Fatalf("expected an error")
			}
			if want := tc.state; got != want {
				t.Errorf("State() = %s, want %s", got, want)
			}
		})
	}
}

func TestMain(m *testing.M) {
	readFileFn = readFile // Override leds.readFileFn for testing.
	os.Exit(m.Run())
}

var (
	readFileData []byte
	readFileErr  error
)

func readFile(filename string) ([]byte, error) {
	return readFileData, readFileErr
}

func prepareReadFile(data []byte, ok bool) {
	readFileData = data
	if !ok {
		readFileErr = fmt.Errorf("ReadFile error for testing")
	}
	readFileErr = nil
}
