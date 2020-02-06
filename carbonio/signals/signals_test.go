package signals

import (
	"fmt"
	"os"
	"testing"

	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/golib/operators"
)

func TestMain(m *testing.M) {
	// Override function pointers for testing.
	helpers.ReadFileFn = readFile
	helpers.WriteFileFn = writeFile

	os.Exit(m.Run())
}

func TestNewInput(t *testing.T) {
	for _, tc := range []struct {
		desc       string
		ok         bool
		num        uint
		maxNum     uint
		dir        Dir
		padSPI     string
		phantomSPI string
	}{
		// Valid inputs.
		{"signal 1", true,
			1, 16, Input,
			"/sys/bus/spi/devices/spi1.1/ch0_pad_en",
			"/sys/bus/spi/devices/spi4.0/adc1_phantom_en"},
		{"signal 2", true,
			2, 16, Input,
			"/sys/bus/spi/devices/spi1.1/ch1_pad_en",
			"/sys/bus/spi/devices/spi4.0/adc1_phantom_en"},
		{"signal 15", true,
			15, 16, Input,
			"/sys/bus/spi/devices/spi1.2/ch2_pad_en",
			"/sys/bus/spi/devices/spi4.0/adc2_phantom_en"},
		{"signal 16", true,
			16, 16, Input,
			"/sys/bus/spi/devices/spi1.2/ch3_pad_en",
			"/sys/bus/spi/devices/spi4.0/adc2_phantom_en"},
		// Invalid inputs.
		{"Number out of range", false, 99, 16, Input, "", ""},
		{"MaxNumber not set", false, 16, 0, Input, "", ""},
	} {
		t.Run(fmt.Sprintf("New() %s", tc.desc), func(t *testing.T) {
			s, err := New("Beep-ba-beep",
				Number(tc.num),
				MaxNumber(tc.maxNum),
				Direction(tc.dir),
			)
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatal("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := s.padSPI, tc.padSPI; got != want {
				t.Errorf("padSPI = %s, want %s", got, want)
			}
			if got, want := s.phantomSPI, tc.phantomSPI; got != want {
				t.Errorf("phantomSPI = %s, want %s", got, want)
			}
		})
	}
}

func TestPad(t *testing.T) {
	for _, tc := range []struct {
		desc      string
		ok        bool
		input     []byte
		isEnabled bool
	}{
		// Supported states.
		{"off", true, []byte{'0', '\n'}, false},
		{"on", true, []byte{'1', '\n'}, true},

		// Error states.
		{"unsupported data", false, []byte{0xff, '\n'}, false},
		{"empty file", false, []byte{}, false},
	} {
		signal, err := New("TestPad",
			Number(1),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("Pad() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.input, tc.ok)
			got, err := signal.Pad()
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if want := tc.isEnabled; got != want {
				t.Errorf("= %t, want %t", got, want)
			}
		})
	}
}

func TestSetPad(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		ok     bool
		enable bool
		output []byte
	}{
		// Supported states.
		{"off", true, false, []byte{'0', '\n'}},
		{"on", true, true, []byte{'1', '\n'}},

		// Error states.
		{desc: "writefile error", ok: false, enable: false},
	} {
		signal, err := New("TestSetPad",
			Number(1),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("SetPad() %s", tc.desc), func(t *testing.T) {
			prepareWriteFile(tc.ok)
			err := signal.SetPad(tc.enable)
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := wfData, tc.output; !operators.EqualSlicesOfByte(got, want) {
				t.Errorf("expected %v to be written, not %v", want, got)
			}
		})
	}
}

func TestPhantom(t *testing.T) {
	for _, tc := range []struct {
		desc      string
		ok        bool
		num       uint
		input     []byte
		isEnabled bool
	}{
		// Supported states. "only" indicates that only that phantom is on/off, and
		// "all" indicates that "all" phantoms for the `agc` device are in that state.

		{"only 1 on", true, 1, []byte{'8', '\n'}, true},
		{"only 1 off", true, 1, []byte{'7', '\n'}, false},
		{"all on 1 on", true, 1, []byte{'1', '5', '\n'}, true},
		{"all off 1 off", true, 1, []byte{'0', '\n'}, false},

		{"only 2 on", true, 2, []byte{'4', '\n'}, true},
		{"only 2 off", true, 2, []byte{'1', '1', '\n'}, false},
		{"all on 2 on", true, 2, []byte{'1', '5', '\n'}, true},
		{"all off 2 off", true, 2, []byte{'0', '\n'}, false},

		{"only 15 on", true, 15, []byte{'2', '\n'}, true},
		{"only 15 off", true, 15, []byte{'1', '3', '\n'}, false},
		{"all on 15 on", true, 15, []byte{'1', '5', '\n'}, true},
		{"all off 15 off", true, 15, []byte{'0', '\n'}, false},

		{"only 16 on", true, 16, []byte{'1', '\n'}, true},
		{"only 16 off", true, 16, []byte{'1', '4', '\n'}, false},
		{"all on 16 on", true, 16, []byte{'1', '5', '\n'}, true},
		{"all off 16 off", true, 16, []byte{'0', '\n'}, false},

		// Error states.
		{"unsupported data", false, 1, []byte{0xff, '\n'}, false},
		{"empty file", false, 1, []byte{}, false},
	} {
		signal, err := New("TestPhantom",
			Number(tc.num),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("Phantom() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.input, tc.ok)
			got, err := signal.Phantom()
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if want := tc.isEnabled; got != want {
				t.Errorf("= %t, want %t", got, want)
			}
		})
	}
}

func TestSetPhantom(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		ok     bool
		num    uint
		input  []byte
		enable bool
		output []byte
	}{
		// Supported states. "none" or "all" indicate that none or all four phantoms
		// are enabled on the given `agc` device.

		{"none 1 on", true, 1, []byte{'0', '\n'}, true, []byte{'8', '\n'}},
		{"none 1 off", true, 1, []byte{'0', '\n'}, false, []byte{'0', '\n'}},
		{"all 1 on", true, 1, []byte{'1', '5', '\n'}, true, []byte{'1', '5', '\n'}},
		{"all 1 off", true, 1, []byte{'1', '5', '\n'}, false, []byte{'7', '\n'}},

		{"none 2 on", true, 2, []byte{'0', '\n'}, true, []byte{'4', '\n'}},
		{"none 2 off", true, 2, []byte{'0', '\n'}, false, []byte{'0', '\n'}},
		{"all 2 on", true, 2, []byte{'1', '5', '\n'}, true, []byte{'1', '5', '\n'}},
		{"all 2 off", true, 2, []byte{'1', '5', '\n'}, false, []byte{'1', '1', '\n'}},

		{"none 15 on", true, 15, []byte{'0', '\n'}, true, []byte{'2', '\n'}},
		{"none 15 off", true, 15, []byte{'0', '\n'}, false, []byte{'0', '\n'}},
		{"all 15 on", true, 15, []byte{'1', '5', '\n'}, true, []byte{'1', '5', '\n'}},
		{"all 15 off", true, 15, []byte{'1', '5', '\n'}, false, []byte{'1', '3', '\n'}},

		{"none 16 on", true, 16, []byte{'0', '\n'}, true, []byte{'1', '\n'}},
		{"none 16 off", true, 16, []byte{'0', '\n'}, false, []byte{'0', '\n'}},
		{"all 16 on", true, 16, []byte{'1', '5', '\n'}, true, []byte{'1', '5', '\n'}},
		{"all 16 off", true, 16, []byte{'1', '5', '\n'}, false, []byte{'1', '4', '\n'}},

		// Error states.
		{"unsupported data", false, 1, []byte{0xff, '\n'}, false, []byte{0, '\n'}},
		{"empty file", false, 1, []byte{}, false, []byte{}},
	} {
		signal, err := New("TestSetPhantom",
			Number(tc.num),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("SetPhantom() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.input, tc.ok)
			prepareWriteFile(tc.ok)
			err := signal.SetPhantom(tc.enable)
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := wfData, tc.output; !operators.EqualSlicesOfByte(got, want) {
				t.Errorf("expected %v to be written, not %v", want, got)
			}
		})
	}
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
	wfErr = nil
	if !ok {
		wfErr = fmt.Errorf("WriteFile error for testing")
	}
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
