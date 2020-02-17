package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/avid-s3l/carbonio/leds"
	"github.com/kward/avid-s3l/carbonio/signals"
	"github.com/spf13/cobra"
)

func init() {
	internalCmd := &cobra.Command{
		Use:    "internal",
		Short:  "internal commands",
		Long:   `internal commands that are used for testing.`,
		Hidden: true,
	}
	rootCmd.AddCommand(internalCmd)

	internalCmd.AddCommand(&cobra.Command{
		Use:   "create_spi",
		Short: "create_setup directory structure for testing",
		Long:  "Create a full spi directory structure for testing.",
		Run:   internal_create_spi,
	})
}

type data struct {
	path string
	data string
}

func internal_create_spi(cmd *cobra.Command, args []string) {
	if spiBaseDir == devices.SPIDevicesDir {
		exit(fmt.Sprintf("refusing to overwrite core SPI dir %s", spiBaseDir))
	}

	tree := []data{}

	ls, err := leds.New(leds.SPIBaseDir(spiBaseDir))
	if err != nil {
		exit(fmt.Sprintf("error instantiating leds; %s", err))
	}
	tree = append(tree, data{ls.Power().Path(), "0"})
	tree = append(tree, data{ls.Status().Path(), "0"})
	tree = append(tree, data{ls.Mute().Path(), "0"})

	for i := uint(1); i <= device.NumMicInputs(); i++ {
		tree = append(tree, data{signals.GainPath(i), "1"})
		tree = append(tree, data{signals.PadPath(i), "0"})
		tree = append(tree, data{signals.PhantomPath(i), "0"})
	}

	fmt.Printf("spiBaseDir: %s\n", spiBaseDir)
	for _, t := range tree {
		fmt.Printf("path: %s data: %q\n", t.path, t.data)
		p := filepath.Join(spiBaseDir, t.path)
		if !dryRun {
			d := path.Dir(p)
			if err := os.MkdirAll(path.Dir(p), 0755); err != nil {
				exit(fmt.Sprintf("error creating directory %s; %v", d, err))
			}
			helpers.WriteSPIFile(p, t.data)
		}
	}
}
