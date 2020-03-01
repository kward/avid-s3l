package signals

import (
	"fmt"
	"testing"

	"github.com/kward/avid-s3l/carbonio/helpers"
)

func TestGain(t *testing.T) {
	signal, err := newInput("TestGain", 1, 16)
	if err != nil {
		t.Fatalf("error setting up test; %s", err)
	}

	for _, tc := range []struct {
		desc     string
		ok       bool
		rfErr    error // MockReadFile error.
		spiValue int   // Current SPI value.

		gain uint
	}{
		// Supported values.
		{"13 dB", true, nil, 4, 13},
		{"21 dB", true, nil, 12, 21},
		{"34 dB", true, nil, 25, 34},
		{"55 dB", true, nil, 46, 55},

		// Error states.
		{desc: "spi value too low", spiValue: 0},
		{desc: "spi value too high", spiValue: 52},
		{desc: "readfile error", rfErr: fmt.Errorf("mock ReadFile error")},
	} {
		t.Run(fmt.Sprintf("Gain() %s", tc.desc), func(t *testing.T) {
			helpers.ResetMockReadWrite()
			helpers.PrepareMockReadFile([]byte{}, tc.rfErr)
			signal.Pad().spi.Write(tc.spiValue)

			got, err := signal.Gain().Value()
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
	signal, err := newInput("TestGain", 1, 16)
	if err != nil {
		t.Fatalf("error setting up test; %s", err)
	}

	for _, tc := range []struct {
		desc  string
		ok    bool
		wfErr error // MockWriteFile error.

		gain  uint
		value int
	}{
		// Supported values.
		{"11 dB", true, nil, 11, 2},
		{"13 dB", true, nil, 13, 4},
		{"17 dB", true, nil, 17, 8},
		//...
		{"47 dB", true, nil, 47, 38},
		{"53 dB", true, nil, 53, 44},
		{"59 dB", true, nil, 59, 50},

		// Error states.
		{desc: "7 dB is too low", gain: 7},
		{desc: "61 dB is too high", gain: 61},
		{desc: "writefile error", wfErr: fmt.Errorf("mock WriteFile error")},
	} {
		t.Run(fmt.Sprintf("Gain() %s", tc.desc), func(t *testing.T) {
			helpers.ResetMockReadWrite()
			// TODO(2020-02-23) Stop writing files directly.
			helpers.PrepareMockWriteFile(tc.wfErr)
			err := signal.Gain().SetValue(tc.gain)
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := signal.Gain().spi.Value(), tc.value; got != want {
				t.Errorf("SPI Value() = %d, want %d", want, got)
			}
		})
	}
}

func TestPad(t *testing.T) {
	signal, err := newInput("TestPad", 1, 16)
	if err != nil {
		t.Fatalf("error setting up test; %s", err)
	}

	for _, tc := range []struct {
		desc     string
		ok       bool
		rfErr    error // MockReadFile error.
		spiValue int   // Current SPI value.

		isEnabled bool
	}{
		// Supported states.
		{"off", true, nil, 0, false},
		{"on", true, nil, 1, true},

		// Error states.
		{desc: "unsupported spi value", spiValue: 123},
		{desc: "readfile error", rfErr: fmt.Errorf("mock ReadFile error")},
	} {
		t.Run(fmt.Sprintf("Pad() %s", tc.desc), func(t *testing.T) {
			helpers.ResetMockReadWrite()
			helpers.PrepareMockReadFile([]byte{}, tc.rfErr)
			signal.Pad().spi.Write(tc.spiValue)

			got, err := signal.Pad().IsEnabled()
			if err != nil && tc.ok {
				t.Fatalf("unexpected error; %s", err)
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

func TestPad_SetPad(t *testing.T) {
	signal, err := newInput("TestSetPad", 1, 16)
	if err != nil {
		t.Fatalf("error setting up test; %s", err)
	}

	for _, tc := range []struct {
		desc  string
		ok    bool
		wfErr error // MockWriteFile error.

		enable bool
		value  int
	}{
		// Supported states.
		{"off", true, nil, PadDisabled, 0},
		{"on", true, nil, PadEnabled, 1},

		// Error states.
		{desc: "writefile error", wfErr: fmt.Errorf("mock WriteFile error")},
	} {
		t.Run(fmt.Sprintf("SetPad() %s", tc.desc), func(t *testing.T) {
			helpers.ResetMockReadWrite()
			helpers.PrepareMockWriteFile(tc.wfErr)

			// Calling setState() directly as [En|Dis]able are simple enough.
			err := signal.Pad().setState(tc.enable)
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := signal.Pad().spi.Value(), tc.value; got != want {
				t.Errorf("SPI Value() = %d, want %d", want, got)
			}
		})
	}
}

func TestPhantom_IsEnabled(t *testing.T) {
	for _, tc := range []struct {
		desc     string
		ok       bool
		rfErr    error // MockReadFile error.
		spiValue int   // Current SPI value.

		num       int
		isEnabled bool
	}{
		// Supported states. "only" indicates that only that phantom is on/off, and
		// "all" indicates that "all" phantoms for the `agc` device are in that state.
		{"only 1 on", true, nil, 0b00001000, 1, true},
		{"only 1 off", true, nil, 0b00000111, 1, false},
		{"all on 1 on", true, nil, 0b00001111, 1, true},
		{"all off 1 off", true, nil, 0b00000000, 1, false},

		{"only 2 on", true, nil, 0b00000100, 2, true},
		{"only 2 off", true, nil, 0b00001011, 2, false},
		{"all on 2 on", true, nil, 0b00001111, 2, true},
		{"all off 2 off", true, nil, 0b00000000, 2, false},

		{"only 15 on", true, nil, 0b00000010, 15, true},
		{"only 15 off", true, nil, 0b00001101, 15, false},
		{"all on 15 on", true, nil, 0b00001111, 15, true},
		{"all off 15 off", true, nil, 0b00000000, 15, false},

		{"only 16 on", true, nil, 0b00000001, 16, true},
		{"only 16 off", true, nil, 0b00001110, 16, false},
		{"all on 16 on", true, nil, 0b00001111, 16, true},
		{"all off 16 off", true, nil, 0b00000000, 16, false},

		// Error states.
		{desc: "unsupported spi value", spiValue: 99, num: 1},
		{desc: "readfile error", rfErr: fmt.Errorf("mock ReadFile error"), num: 1},
	} {
		signal, err := newInput("TestPhantom", tc.num, 16)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("Phantom() %s", tc.desc), func(t *testing.T) {
			helpers.ResetMockReadWrite()
			helpers.PrepareMockReadFile([]byte{}, tc.rfErr)
			signal.Phantom().spi.Write(tc.spiValue)

			got, err := signal.Phantom().IsEnabled()
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
		rfErr    error // MockReadFile error.
		wfErr    error // MockWriteFile error.
		spiValue int   // Current SPI value.

		num    int
		enable bool
		value  int // New SPI value.
	}{
		// Supported states. "none" or "all" indicate that none or all four phantoms
		// are enabled on the given `agc` device.
		{"none 1 on", true, nil, nil, 0b00000000, 1, true, 0b00001000},
		{"none 1 off", true, nil, nil, 0b00000000, 1, false, 0b00000000},
		{"all 1 on", true, nil, nil, 0b00001111, 1, true, 0b00001111},
		{"all 1 off", true, nil, nil, 0b00001111, 1, false, 0b00000111},

		{"none 2 on", true, nil, nil, 0b00000000, 2, true, 0b00000100},
		{"none 2 off", true, nil, nil, 0b00000000, 2, false, 0b00000000},
		{"all 2 on", true, nil, nil, 0b00001111, 2, true, 0b00001111},
		{"all 2 off", true, nil, nil, 0b00001111, 2, false, 0b00001011},

		{"none 15 on", true, nil, nil, 0b00000000, 15, true, 0b00000010},
		{"none 15 off", true, nil, nil, 0b00000000, 15, false, 0b00000000},
		{"all 15 on", true, nil, nil, 0b00001111, 15, true, 0b00001111},
		{"all 15 off", true, nil, nil, 0b00001111, 15, false, 0b00001101},

		{"none 16 on", true, nil, nil, 0b00000000, 16, true, 0b00000001},
		{"none 16 off", true, nil, nil, 0b00000000, 16, false, 0b00000000},
		{"all 16 on", true, nil, nil, 0b00001111, 16, true, 0b00001111},
		{"all 16 off", true, nil, nil, 0b00001111, 16, false, 0b00001110},

		// Error states. A ReadFile error is included because the phantom value
		// must be read before it can be written.
		{desc: "readfile error", rfErr: fmt.Errorf("mock ReadFile error"), num: 1},
		{desc: "writefile error", wfErr: fmt.Errorf("mock WriteFile error"), num: 1},
	} {
		signal, err := newInput("TestSetPhantom", tc.num, 16)
		if err != nil {
			t.Fatalf("error setting up test; %s", err)
		}

		t.Run(fmt.Sprintf("SetPhantom() %s", tc.desc), func(t *testing.T) {
			helpers.ResetMockReadWrite()
			helpers.PrepareMockReadFile([]byte{}, tc.rfErr)
			helpers.PrepareMockWriteFile(tc.wfErr)
			signal.Phantom().spi.Write(tc.spiValue)

			// Calling setState() directly as [En|Dis]able are simple enough.
			err := signal.Phantom().setState(tc.enable)
			if err != nil && tc.ok {
				t.Fatalf("unexpected error %q", err)
			}
			if err == nil && !tc.ok {
				t.Fatalf("expected an error")
			}
			if !tc.ok {
				return
			}
			if got, want := signal.Phantom().spi.Value(), tc.value; got != want {
				t.Errorf("SPI Value() = %d, want %d", want, got)
			}
		})
	}
}
