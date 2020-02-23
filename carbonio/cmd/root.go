package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/avid-s3l/carbonio/helpers"
	"github.com/kward/avid-s3l/carbonio/spi"
	"github.com/spf13/cobra"
)

var (
	dryRun     bool
	spiBaseDir string
	verbose    bool

	device devices.Device

	rootCmd = &cobra.Command{
		Use:   "carbonio",
		Short: "carbonio provides direct control of the Avid Carbon I/O device.",
		Long: `carbonio provides direct control of the Avid Carbon I/O device, which is built
into the E3 Engine and the Stage 16 stage boxes. Complete documentation is available at
https://github.com/kward/avid-s3l`,
		PersistentPreRun: persistentPreRun,
		Run:              root,
	}
)

const (
	ifs = " "
	ofs = " "
)

func Execute() {
	rootCmd.PersistentFlags().BoolVarP(
		&dryRun, "dry_run", "n", false, "perform a dry-run")
	rootCmd.PersistentFlags().BoolVarP(
		&verbose, "verbose", "v", false, "verbose output")

	// TODO(2020-02-18) If this flag is overridden, validate a "version" file to
	// ensure the structure is appropriate for testing.
	rootCmd.PersistentFlags().StringVarP(
		&spiBaseDir, "spi_base_dir", "", spi.DevicesDir, "spi base directory")

	if err := rootCmd.Execute(); err != nil {
		helpers.Exit(fmt.Sprintf("error: %v", err))
	}
}

func persistentPreRun(cmd *cobra.Command, args []string) {
	var err error

	if cmd.Use == "help" || cmd.Use == "help [command]" {
		return
	}

	// Setup carbonio device.
	// NOTE: declaring with '=' to ensure global `device` is not overridden.
	spiDelayRead := false
	if cmd.Use == "create_spi" {
		spiDelayRead = true
	}
	device, err = devices.NewStage16(
		devices.SPIDelayRead(spiDelayRead),
		devices.SPIBaseDir(spiBaseDir),
		devices.Verbose(verbose),
	)
	if err != nil {
		helpers.Exit(fmt.Sprintf("error configuring the Stage 16 device; %s", err))
	}

	// Validate spi_base_dir.
	if cmd.Use == "create_spi" { // We're going to create the dir.
		return
	}
	err = filepath.Walk(spiBaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		helpers.Exit(fmt.Sprintf("invalid --spi_base_dir flag value %s", spiBaseDir))
	}

}

func root(cmd *cobra.Command, args []string) {
	fmt.Println("This is the carbonio controller.")
	if device == nil {
		fmt.Println("device is unitialized")
		return
	}
}
