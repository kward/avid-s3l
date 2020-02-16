package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/helpers"
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
		Use:   "setup_spi",
		Short: "setup spi directory structure for testing",
		Long:  "Create a full spi directory structure for testing.",
		Run:   internal_setup_spi,
	})
}

type data struct {
	path string
	data string
}

func internal_setup_spi(cmd *cobra.Command, args []string) {
	if spiBaseDir == devices.SPIDevicesDir {
		exit(fmt.Sprintf("refusing to overwrite core SPI dir %s", spiBaseDir))
	}

	tree := []data{}
	for i := uint(1); i <= device.NumMicInputs(); i++ {
		tree = append(tree, data{
			path: signals.PadPath(i),
			data: "0",
		})
	}

	for _, t := range tree {
		fmt.Printf("spiBaseDir: %s path: %s data: %q\n", spiBaseDir, t.path, t.data)
		p := filepath.Join(spiBaseDir, t.path)
		if !dryRun {
			d := path.Dir(p)
			if err := os.MkdirAll(path.Dir(p), 0755); err != nil {
				exit(fmt.Sprintf("error creating directory %s; %v", d, err))
			}
			helpers.WriteFile(p, t.data)
		}
	}
}
