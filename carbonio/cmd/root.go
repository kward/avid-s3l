package cmd

import (
	"fmt"
	"os"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/spf13/cobra"
)

var (
	spiBaseDir string
	verbose    bool

	rootCmd = &cobra.Command{
		Use:   "carbonio",
		Short: "carbonio provides direct control of the Avid Carbon I/O device.",
		Long: `carbonio provides direct control of the Avid Carbon I/O device, which is built
into the E3 Engine and the Stage 16 stage boxes. Complete documentation is available at
https://github.com/kward/avid-s3l`,
		Run: root,
	}
)

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false, "verbose output")

	rootCmd.PersistentFlags().StringVarP(
		&spiBaseDir, "spi_base_dir", "", devices.SPIDevicesDir, "spi base directory")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func root(cmd *cobra.Command, args []string) {
	fmt.Println("This is the carbonio controller.")

	// // LEDs.
	// ls := []*leds.LED{leds.Power, leds.Status, leds.Mute}
	// for _, l := range ls {
	//  fmt.Println(l)
	// }
	// fmt.Println("Toggling LEDs…")
	// for _, l := range ls {
	//  l.SetState(leds.Off)
	// }
	// time.Sleep(2 * time.Second)
	// for _, l := range ls {
	//  l.SetState(leds.On)
	// }
}
