/*
Package helpers provides code snippets that are used across the code base.
*/
package helpers

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/kward/golib/math"
)

const maxDataLen = 32

var (
	fileMode = os.FileMode(0644)
	// readFileFn holds a pointer to an ioutil.ReadFile compatible function.
	readFileFn readFileType
	// WriteFileFn holds a pointer to an ioutil.WriteFile compatible function.
	writeFileFn writeFileType
)

func init() {
	SetReadFileFn(ioutil.ReadFile)
	SetWriteFileFn(ioutil.WriteFile)
}

type readFileType func(filename string) ([]byte, error)
type writeFileType func(filename string, data []byte, perm os.FileMode) error

func SetReadFileFn(fn readFileType)   { readFileFn = fn }
func SetWriteFileFn(fn writeFileType) { writeFileFn = fn }

// ReadByte from the SPI interface.
func ReadByte(filename string) (byte, error) {
	data, err := readFileFn(filename)
	if err != nil {
		return 0, err
	}
	if len(data) != 2 || data[1] != '\n' {
		return 0, fmt.Errorf("%q contains unexpected data; %v", filename, data[0:math.MinInt(len(data), maxDataLen)])
	}
	return data[0], nil
}

// ReadSPIFile reads data from the SPI interface, with newline stripped.
func ReadSPIFile(filename string) (string, error) {
	data, err := readFileFn(filename)
	if err != nil {
		return "", err
	}
	return strings.Split(fmt.Sprintf("%s", data), "\n")[0], nil
}

// WriteByte to the SPI interface.
func WriteByte(filename string, v byte) error {
	return writeFileFn(filename, []byte{v, '\n'}, fileMode)
}

// WriteSPIFile writes data to the SPI interface.
func WriteSPIFile(filename string, data string) error {
	buf := bytes.NewBuffer(nil)
	_, err := buf.WriteString(data)
	if err != nil {
		return fmt.Errorf("error writing string to buffer; %s", err)
	}
	err = buf.WriteByte('\n')
	if err != nil {
		return fmt.Errorf("error writing newline to buffer; %s", err)
	}
	return writeFileFn(filename, buf.Bytes(), fileMode)
}

//-----------------------------------------------------------------------------
// Helpers for testing.

var (
	rfData, wfData []byte
	rfErr, wfErr   error
)

func PrepareReadFile(data []byte, err error) {
	if err != nil {
		rfErr = err
		return
	}
	rfData = data
	rfErr = nil
}

// MockReadFile matches the signature of io.ReadFile.
func MockReadFile(filename string) ([]byte, error) {
	return rfData, rfErr
}

func PrepareWriteFile(err error) {
	if err != nil {
		wfErr = err
		return
	}
	wfErr = nil
}

// MockWriteFile matches the signature of io.WriteFile.
func MockWriteFile(filename string, data []byte, mode os.FileMode) error {
	if wfErr != nil {
		wfData = []byte{}
		return wfErr
	}
	wfData = data
	return nil
}

func MockWriteData() []byte { return wfData }
