/*
Package helpers provides code snippets that are used across the code base.
*/
package helpers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// CommonLogFormat generates a log entry in the Common Log Format (CLF).
//   "%h %l %u %t \"%r\" %>s %b"
func CommonLogFormat(r *http.Request, status, length int) {
	log.Printf("%s - - %s %q %d %d\n",
		r.RemoteAddr, time.Now().Format("02/Jan/2006:15:04:05 -0700"), r.URL.Path, status, length)
}

func Exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

//-----------------------------------------------------------------------------
// Helpers for testing.

func init() {
	SetReadFileFn(ioutil.ReadFile)
	SetWriteFileFn(ioutil.WriteFile)
}

var (
	// readFileFn holds a pointer to an ioutil.ReadFile compatible function.
	readFileFn readFileType
	// WriteFileFn holds a pointer to an ioutil.WriteFile compatible function.
	writeFileFn writeFileType
)

type readFileType func(filename string) ([]byte, error)
type writeFileType func(filename string, data []byte, perm os.FileMode) error

func ReadFileFn() readFileType   { return readFileFn }
func WriteFileFn() writeFileType { return writeFileFn }

// SetReadFileFn overrides the ReadFile function.
func SetReadFileFn(fn readFileType) { readFileFn = fn }

// SetWriteFileFn overrides the WriteFile function.
func SetWriteFileFn(fn writeFileType) { writeFileFn = fn }

var (
	rfData       []byte
	rfErr, wfErr error
)

// PrepareMockReadFile to return error for ReadFileFn() call.
func PrepareMockReadFile(data []byte, err error) {
	if err != nil {
		rfErr = err
		return
	}
	rfData = data
	rfErr = nil
}

// MockReadFile matches the signature of io.ReadFile.
func MockReadFile(filename string) ([]byte, error) { return rfData, rfErr }

// MockData returns the mock ReadFileFn() data.
func MockData() []byte { return rfData }

// PrepareMockWriteFile to return error for WriteFileFn() call.
func PrepareMockWriteFile(err error) {
	if err != nil {
		wfErr = err
		return
	}
	wfErr = nil
}

// MockWriteFile matches the signature of io.WriteFile.
func MockWriteFile(filename string, data []byte, mode os.FileMode) error {
	if wfErr != nil {
		rfData = []byte{}
		return wfErr
	}
	rfData = data
	return nil
}

// ResetMockReadWrite to clean state.
func ResetMockReadWrite() {
	rfData = []byte{}
	rfErr = nil
	wfErr = nil
}
