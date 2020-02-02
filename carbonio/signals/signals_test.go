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
		desc    string
		num     int
		dir     Dir
		iface   string
		channel string
		ok      bool
	}{
		// Valid inputs.
		{"signal 1", 1, Input, "/sys/bus/spi/devices/spi1.1", "ch0", true},
		{"signal 16", 16, Input, "/sys/bus/spi/devices/spi1.2", "ch3", true},
		// Invalid inputs.
		{"Number out of range", 0, Input, "", "", false},
		{"MaxNumber not set", 99, Input, "", "", false},
	} {
		t.Run(fmt.Sprintf("New() %s", tc.desc), func(t *testing.T) {
			s, err := New("Beep-ba-beep",
				Number(tc.num),
				MaxNumber(16),
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
			if got, want := s.iface, tc.iface; got != want {
				t.Errorf("iface = %s, want %s", got, want)
			}
			if got, want := s.channel, tc.channel; got != want {
				t.Errorf("channel = %s, want %s", got, want)
			}
		})
	}
}

func TestPad(t *testing.T) {
	for _, tc := range []struct {
		desc string
		pad  bool
		data []byte
		ok   bool
	}{
		// Supported states.
		{"off", false, []byte{'0', '\n'}, true},
		{"on", true, []byte{'1', '\n'}, true},
		// TODO: The following are only really useful for ReadFile as the local
		// WriteFile doesn't validate what should have been written.
		{"unsupported", false, []byte{123, '\n'}, false},
		{"readfile error", false, []byte{}, false},
	} {
		signal, err := New("Pad test",
			Number(1),
			MaxNumber(16),
			Direction(Input),
		)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("Pad() %s", tc.desc), func(t *testing.T) {
			prepareReadFile(tc.data, tc.ok)
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
			if want := tc.pad; got != want {
				t.Errorf("= %t, want %t", got, want)
			}
		})

		t.Run(fmt.Sprintf("SetPad() %s", tc.desc), func(t *testing.T) {
			prepareWriteFile(tc.ok)
			err := signal.SetPad(tc.pad)
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
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
