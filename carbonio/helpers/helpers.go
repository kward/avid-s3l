/*
Package helpers provides code snippets that are used across the code base.
*/
package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/kward/golib/math"
)

const maxDataLen = 32

var fileMode = os.FileMode(0644)

// ReadFileFn holds a pointer to the ioutil.ReadFile function. This pointer can
// be overridden for testing.
var ReadFileFn = ioutil.ReadFile

// ReadByte from the SPI interface.
func ReadByte(filename string) (byte, error) {
	data, err := ReadFileFn(filename)
	if err != nil {
		return 0, err
	}
	if len(data) != 2 || data[1] != '\n' {
		return 0, fmt.Errorf("%q contains unexpected data; %v", filename, data[0:math.MinInt(len(data), maxDataLen)])
	}
	return data[0], nil
}

// WriteFileFn holds a pointer to the ioutil.WriteFile function. This pointer
// can be overridden for testing.
var WriteFileFn = ioutil.WriteFile

// WriteByte to the SPI interface.
func WriteByte(filename string, v byte) error {
	return WriteFileFn(filename, []byte{v, '\n'}, fileMode)
}

// WriteFile to the SPI interface.
func WriteFile(filename string, data string) error {
	buf := bytes.NewBuffer(nil)
	_, err := buf.WriteString(data)
	if err != nil {
		return fmt.Errorf("error writing string to buffer; %s", err)
	}
	err = buf.WriteByte('\n')
	if err != nil {
		return fmt.Errorf("error writing newline to buffer; %s", err)
	}
	return WriteFileFn(filename, buf.Bytes(), fileMode)
}
