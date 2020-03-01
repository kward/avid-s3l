package signals

import (
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
