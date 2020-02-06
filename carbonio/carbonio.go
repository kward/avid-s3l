// The carbionio command enables control of Carbon I/O hardware.
//
// Build for ARM with
// $ GOOS=linux GOARM=7 GOARCH=arm go build carbonio.go
package main

import (
	"github.com/kward/avid-s3l/carbonio/cmd"
)

func main() {
	cmd.Execute()
}
