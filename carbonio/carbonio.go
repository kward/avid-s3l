// The carbionio command enables control of the Avid Stage 16 device I/O and
// LEDs.
package main

import (
	"fmt"

	"github.com/kward/avid-s3l/carbonio/carbonio/leds"
)

func main() {
	fmt.Println("Hello, world!")
	fmt.Printf("LED: %s\n", leds.On)
}
