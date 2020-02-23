package spi

import (
	"fmt"
	"os"
	"testing"

	"github.com/kward/avid-s3l/carbonio/helpers"
)

func TestMain(m *testing.M) {
	// Override function pointers for testing.
	helpers.SetReadFileFn(helpers.MockReadFile)
	helpers.SetWriteFileFn(helpers.MockWriteFile)

	os.Exit(m.Run())
}

func TestSPI(t *testing.T) {
	for _, tc := range []struct {
		desc   string
		ok     bool
		rfData []byte // MockReadFile data.
		rfErr  error  // MockReadFile error.

		enum  Enum
		num   int
		path  string
		value int
	}{
		// Valid values.
		{"gain #1", true, []byte("1\n"), nil,
			Gain, 1, "/spi/base/spi1.1/ch0_preamp_gain", 1},
		{"gain #2", true, []byte("2\n"), nil,
			Gain, 2, "/spi/base/spi1.1/ch1_preamp_gain", 2},
		{"gain #15", true, []byte("50\n"), nil,
			Gain, 15, "/spi/base/spi1.2/ch2_preamp_gain", 50},
		{"gain #16", true, []byte("51\n"), nil,
			Gain, 16, "/spi/base/spi1.2/ch3_preamp_gain", 51},

		{"pad #1", true, []byte("0\n"), nil,
			Pad, 1, "/spi/base/spi1.1/ch0_pad_en", 0},
		{"pad #2", true, []byte("1\n"), nil,
			Pad, 2, "/spi/base/spi1.1/ch1_pad_en", 1},
		{"pad #15", true, []byte("0\n"), nil,
			Pad, 15, "/spi/base/spi1.2/ch2_pad_en", 0},
		{"pad #16", true, []byte("1\n"), nil,
			Pad, 16, "/spi/base/spi1.2/ch3_pad_en", 1},

		{"phantom #1", true, []byte("0\n"), nil,
			Phantom, 1, "/spi/base/spi4.0/adc1_phantom_en", 0},
		{"phantom #2", true, []byte("1\n"), nil,
			Phantom, 2, "/spi/base/spi4.0/adc1_phantom_en", 1},
		{"phantom #15", true, []byte("0\n"), nil,
			Phantom, 15, "/spi/base/spi4.0/adc2_phantom_en", 0},
		{"phantom #16", true, []byte("1\n"), nil,
			Phantom, 16, "/spi/base/spi4.0/adc2_phantom_en", 1},

		// Invalid values.
		// TODO(2020-02-22) Add some invalid values from signals_test.
		{desc: "unsupported data", rfData: []byte{0xff}},
		{desc: "too much data", rfData: []byte("18446744073709551615\n")}, // Max uint64.
		{desc: "wrong termination", rfData: []byte{'2', 34}},
		{desc: "empty file"},
		{desc: "readfile error", rfErr: fmt.Errorf("mock ReadFile error")},
	} {
		var s *SPI

		helpers.ResetMockReadWrite()
		helpers.PrepareMockReadFile(tc.rfData, tc.rfErr)
		t.Run(fmt.Sprintf("New() %s", tc.desc), func(t *testing.T) {
			s, err := New(tc.enum, tc.num,
				DelayRead(true),
				BaseDir("/spi/base"),
			)
			if err != nil {
				t.Fatalf("New() unexpected error; %s", err)
			}
			if !tc.ok {
				return
			}
			if got, want := s.path, tc.path; got != want {
				t.Errorf("s.path = %s, want %s", got, want)
			}
		})
		if s == nil {
			continue
		}

		// Technically, Read() is called from within New(), but we'll test it more.
		t.Run(fmt.Sprintf("Read() %s", tc.desc), func(t *testing.T) {
			v, err := s.Read()
			if err != nil && tc.ok {
				t.Fatalf("unexpected error; %s", err)
			}
			if err == nil && !tc.ok {
				t.Fatal("expected an error")
			}
			if !tc.ok {
				return
			}
			fmt.Printf("desc: %s, v: %d\n", tc.desc, v)
			if got, want := v, tc.value; got != want {
				t.Errorf("= %d, want %d", got, want)
			}
		})

		t.Run(fmt.Sprintf("Name() %s", tc.desc), func(t *testing.T) {
			if got, want := s.Name(), tc.enum.String(); got != want {
				t.Errorf("= %s, want %s", got, want)
			}
		})
	}
}
