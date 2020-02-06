package cmd

import (
	"fmt"
	"os"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "list the carbonio device settings",
	Long:  `List prints the current settings of the carbonio device.`,
	Run:   list,
}

const ifs = " "

func list(cmd *cobra.Command, args []string) {
	d, err := devices.NewStage16()
	if err != nil {
		fmt.Printf("error configuring the Stage 16 device; %s\n", err)
		os.Exit(1)
	}

	fmt.Println("SIGNAL GAIN PAD PHANTOM")
	for i := uint(1); i <= d.NumMicInputs(); i++ {
		in, err := d.MicInput(i)
		if err != nil {
			fmt.Printf("error accessing mic input %d; %s\n", i, err)
			continue
		}

		gain, err := in.Gain()
		var gainStr string
		if err != nil {
			gainStr = "err"
		} else {
			gainStr = fmt.Sprintf("%d", gain)
		}

		pad, err := in.Pad()
		var padStr string
		if err != nil {
			padStr = "err"
		} else {
			padStr = fmt.Sprintf("%t", pad)
		}

		phantom, err := in.Phantom()
		var phantomStr string
		if err != nil {
			phantomStr = "err"
		} else {
			phantomStr = fmt.Sprintf("%t", phantom)
		}

		fmt.Printf("input/mic/%d %s %s %s\n", i, gainStr, padStr, phantomStr)
	}
}
