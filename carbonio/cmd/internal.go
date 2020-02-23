package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/avid-s3l/carbonio/spi"
	"github.com/spf13/cobra"
)

var (
	internalCmd = &cobra.Command{
		Use:    "internal",
		Short:  "internal commands",
		Long:   `internal commands that are used for testing.`,
		Hidden: true,
	}
)

func init() {
	rootCmd.AddCommand(internalCmd)
	internalCmd.AddCommand(&cobra.Command{
		Use:   "create_spi",
		Short: "create_setup directory structure for testing",
		Long:  "Create a full spi directory structure for testing.",
		Run:   internal_create_spi,
	})
}

func internal_create_spi(cmd *cobra.Command, args []string) {
	if spiBaseDir == spi.DevicesDir {
		helpers.Exit(fmt.Sprintf("refusing to overwrite core SPI dir %s", spiBaseDir))
	}

	// Gather the SPI devices to create.
	type funcs struct {
		name string
		path spi.PathFn
		init spi.InitializeFn
	}
	devs := []funcs{
		{"LED " + device.LEDs().Power().Name(),
			device.LEDs().Power().Path, device.LEDs().Power().Initialize},
		{"LED " + device.LEDs().Status().Name(),
			device.LEDs().Status().Path, device.LEDs().Status().Initialize},
		{"LED " + device.LEDs().Mute().Name(),
			device.LEDs().Mute().Path, device.LEDs().Mute().Initialize},
	}
	for i := 1; i < device.NumMicInputs(); i++ {
		s, err := device.MicInput(i)
		if err != nil {
			fmt.Printf("mic input error; %s", err)
			continue
		}
		devs = append(devs, funcs{fmt.Sprintf("Mic #%d %s", i, s.Gain().Name()),
			s.Gain().Path, s.Gain().Initialize})
		devs = append(devs, funcs{fmt.Sprintf("Mic #%d %s", i, s.Pad().Name()),
			s.Pad().Path, s.Pad().Initialize})
		devs = append(devs, funcs{fmt.Sprintf("Mic #%d %s", i, s.Phantom().Name()),
			s.Phantom().Path, s.Phantom().Initialize})
	}

	// Create the SPI devices.
	fmt.Println("creating…")
	for _, dev := range devs {
		fmt.Printf("  %s: %s\n", dev.name, dev.path())
		if dryRun {
			continue
		}
		dir := path.Dir(dev.path())
		if err := os.MkdirAll(dir, 0755); err != nil {
			helpers.Exit(fmt.Sprintf("error creating directory %s; %v", dir, err))
		}
		dev.init()
	}
}
