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
		gainSPI    string
		padSPI     string
		phantomSPI string
	}{
		// Valid inputs.
		{"signal 1", true,
			1, 16, Input,
			"/spi/base/spi1.1/ch0_preamp_gain",
			"/spi/base/spi1.1/ch0_pad_en",
			"/spi/base/spi4.0/adc1_phantom_en"},
		{"signal 2", true,
			2, 16, Input,
			"/spi/base/spi1.1/ch1_preamp_gain",
			"/spi/base/spi1.1/ch1_pad_en",
			"/spi/base/spi4.0/adc1_phantom_en"},
		{"signal 15", true,
			15, 16, Input,
			"/spi/base/spi1.2/ch2_preamp_gain",
			"/spi/base/spi1.2/ch2_pad_en",
			"/spi/base/spi4.0/adc2_phantom_en"},
		{"signal 16", true,
			16, 16, Input,
			"/spi/base/spi1.2/ch3_preamp_gain",
			"/spi/base/spi1.2/ch3_pad_en",
			"/spi/base/spi4.0/adc2_phantom_en"},

		// Invalid inputs.
		{desc: "Number out of range", num: 99, maxNum: 16, dir: Input},
		{desc: "MaxNumber not set", num: 16, maxNum: 0, dir: Input},
	} {
		t.Run(fmt.Sprintf("New() %s", tc.desc), func(t *testing.T) {
			s, err := New("Beep-ba-beep",
				SPIBaseDir("/spi/base"),
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
			if got, want := s.gainSPI, tc.gainSPI; got != want {
				t.Errorf("gainSPI = %s, want %s", got, want)
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

func TestGain(t *testing.T) {
	for _, tc := range []struct {
		desc    string
		ok      bool
		input   []byte
		readErr error
		gain    uint
	}{
		// Supported values.
		{"13 dB", true, []byte{'4', '\n'}, nil, 13},
		{"21 dB", true, []byte{'1', '2', '\n'}, nil, 21},
		{"34 dB", true, []byte{'2', '5', '\n'}, nil, 34},
		{"55 dB", true, []byte{'4', '6', '\n'}, nil, 55},

		// Error states.
		{desc: "9 dB too low", input: []byte{'0', '\n'}},
		{desc: "89 dB too high", input: []byte{'8', '0', '\n'}},
		{desc: "empty file", input: []byte{}},
		{desc: "readfile error", readErr: fmt.Errorf("some error")},
	} {
		signal, err := New("TestGain",
			SPIBaseDir("/spi/base"),
			Number(1),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("Gain() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.input, tc.readErr)
			got, err := signal.Gain()
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if want := tc.gain; got != want {
				t.Errorf("= %d, want %d", got, want)
			}
		})
	}
}

func TestSetGain(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		ok       bool
		gain     uint
		output   []byte
		writeErr error
	}{
		// Supported values.
		{"11 dB", true, 11, []byte{'2', '\n'}, nil},
		{"13 dB", true, 13, []byte{'4', '\n'}, nil},
		{"17 dB", true, 17, []byte{'8', '\n'}, nil},
		//...
		{"47 dB", true, 47, []byte{'3', '8', '\n'}, nil},
		{"53 dB", true, 53, []byte{'4', '4', '\n'}, nil},
		{"59 dB", true, 59, []byte{'5', '0', '\n'}, nil},

		// Error states.
		{desc: "7 dB is too low"},
		{desc: "61 dB is too high"},
		{desc: "readfile error", writeErr: fmt.Errorf("some error")},
	} {
		signal, err := New("TestGain",
			SPIBaseDir("/spi/base"),
			Number(1),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("Gain() %s", tc.desc), func(t *testing.T) {
			prepareWriteFile(tc.writeErr)
			err := signal.SetGain(tc.gain)
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

func TestPad(t *testing.T) {
	for _, tc := range []struct {
		desc      string
		ok        bool
		input     []byte
		readErr   error
		isEnabled bool
	}{
		// Supported states.
		{"off", true, []byte{'0', '\n'}, nil, false},
		{"on", true, []byte{'1', '\n'}, nil, true},

		// Error states.
		{desc: "unsupported data", input: []byte{0xff, '\n'}},
		{desc: "empty file", input: []byte{}},
		{desc: "readfile error", readErr: fmt.Errorf("ReadFile error")},
	} {
		signal, err := New("TestPad",
			SPIBaseDir("/spi/base"),
			Number(1),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("Pad() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.input, tc.readErr)
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
		desc     string
		ok       bool
		enable   bool
		output   []byte
		writeErr error
	}{
		// Supported states.
		{"off", true, false, []byte{'0', '\n'}, nil},
		{"on", true, true, []byte{'1', '\n'}, nil},

		// Error states.
		{desc: "writefile error", writeErr: fmt.Errorf("WriteFile error")},
	} {
		signal, err := New("TestSetPad",
			SPIBaseDir("/spi/base"),
			Number(1),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("SetPad() %s", tc.desc), func(t *testing.T) {
			prepareWriteFile(tc.writeErr)
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
		readErr   error
		isEnabled bool
	}{
		// Supported states. "only" indicates that only that phantom is on/off, and
		// "all" indicates that "all" phantoms for the `agc` device are in that state.

		{"only 1 on", true, 1, []byte{'8', '\n'}, nil, true},
		{"only 1 off", true, 1, []byte{'7', '\n'}, nil, false},
		{"all on 1 on", true, 1, []byte{'1', '5', '\n'}, nil, true},
		{"all off 1 off", true, 1, []byte{'0', '\n'}, nil, false},

		{"only 2 on", true, 2, []byte{'4', '\n'}, nil, true},
		{"only 2 off", true, 2, []byte{'1', '1', '\n'}, nil, false},
		{"all on 2 on", true, 2, []byte{'1', '5', '\n'}, nil, true},
		{"all off 2 off", true, 2, []byte{'0', '\n'}, nil, false},

		{"only 15 on", true, 15, []byte{'2', '\n'}, nil, true},
		{"only 15 off", true, 15, []byte{'1', '3', '\n'}, nil, false},
		{"all on 15 on", true, 15, []byte{'1', '5', '\n'}, nil, true},
		{"all off 15 off", true, 15, []byte{'0', '\n'}, nil, false},

		{"only 16 on", true, 16, []byte{'1', '\n'}, nil, true},
		{"only 16 off", true, 16, []byte{'1', '4', '\n'}, nil, false},
		{"all on 16 on", true, 16, []byte{'1', '5', '\n'}, nil, true},
		{"all off 16 off", true, 16, []byte{'0', '\n'}, nil, false},

		// Error states.
		{desc: "unsupported data", num: 1, input: []byte{0xff, '\n'}},
		{desc: "unsupported spi value", num: 1, input: []byte{'9', '9', '\n'}},
		{desc: "empty file", num: 1, input: []byte{}},
		{desc: "readfile error", num: 1, readErr: fmt.Errorf("ReadFile error")},
	} {
		signal, err := New("TestPhantom",
			SPIBaseDir("/spi/base"),
			Number(tc.num),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("Phantom() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.input, tc.readErr)
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
		desc     string
		ok       bool
		num      uint
		input    []byte
		readErr  error
		enable   bool
		output   []byte
		writeErr error
	}{
		// Supported states. "none" or "all" indicate that none or all four phantoms
		// are enabled on the given `agc` device.

		{"none 1 on", true, 1, []byte{'0', '\n'}, nil, true, []byte{'8', '\n'}, nil},
		{"none 1 off", true, 1, []byte{'0', '\n'}, nil, false, []byte{'0', '\n'}, nil},
		{"all 1 on", true, 1, []byte{'1', '5', '\n'}, nil, true, []byte{'1', '5', '\n'}, nil},
		{"all 1 off", true, 1, []byte{'1', '5', '\n'}, nil, false, []byte{'7', '\n'}, nil},

		{"none 2 on", true, 2, []byte{'0', '\n'}, nil, true, []byte{'4', '\n'}, nil},
		{"none 2 off", true, 2, []byte{'0', '\n'}, nil, false, []byte{'0', '\n'}, nil},
		{"all 2 on", true, 2, []byte{'1', '5', '\n'}, nil, true, []byte{'1', '5', '\n'}, nil},
		{"all 2 off", true, 2, []byte{'1', '5', '\n'}, nil, false, []byte{'1', '1', '\n'}, nil},

		{"none 15 on", true, 15, []byte{'0', '\n'}, nil, true, []byte{'2', '\n'}, nil},
		{"none 15 off", true, 15, []byte{'0', '\n'}, nil, false, []byte{'0', '\n'}, nil},
		{"all 15 on", true, 15, []byte{'1', '5', '\n'}, nil, true, []byte{'1', '5', '\n'}, nil},
		{"all 15 off", true, 15, []byte{'1', '5', '\n'}, nil, false, []byte{'1', '3', '\n'}, nil},

		{"none 16 on", true, 16, []byte{'0', '\n'}, nil, true, []byte{'1', '\n'}, nil},
		{"none 16 off", true, 16, []byte{'0', '\n'}, nil, false, []byte{'0', '\n'}, nil},
		{"all 16 on", true, 16, []byte{'1', '5', '\n'}, nil, true, []byte{'1', '5', '\n'}, nil},
		{"all 16 off", true, 16, []byte{'1', '5', '\n'}, nil, false, []byte{'1', '4', '\n'}, nil},

		// Error states.
		{desc: "unsupported data", num: 1, input: []byte{0xff, '\n'}},
		{desc: "unsupported spi value", num: 1, input: []byte{'9', '9', '\n'}},
		{desc: "empty file", num: 1, input: []byte{}},
		{desc: "readfile error", num: 1, readErr: fmt.Errorf("ReadFile error")},
		{desc: "readfile error", num: 1, writeErr: fmt.Errorf("WriteFile error")},
	} {
		signal, err := New("TestSetPhantom",
			SPIBaseDir("/spi/base"),
			Number(tc.num),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("SetPhantom() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.input, tc.readErr)
			prepareWriteFile(tc.writeErr)
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

func prepareReadFile(data []byte, err error) {
	rfData = data
	rfErr = err
}

// readFile matches the signature of io.ReadFile.
func readFile(filename string) ([]byte, error) {
	return rfData, rfErr
}

func prepareWriteFile(err error) {
	wfErr = err
}

// writeFile matches the signature of io.WriteFile.
func writeFile(filename string, data []byte, mode os.FileMode) error {
	if wfErr != nil {
		wfData = nil
		return wfErr
	}
	wfData = data
	return nil
}
