package cmd

import (
	"fmt"
	"os"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/kward/tabulate/tabulate"
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

	lines := []string{}
	lines = append(lines, "SIGNAL GAIN PAD PHANTOM")
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

		lines = append(lines, fmt.Sprintf("input/mic/%d %s %s %s", i, gainStr, padStr, phantomStr))
	}
	tbl, err := tabulate.NewTable()
	if err != nil {
		fmt.Printf("unable to list settings; %s", err)
		return
	}
	tbl.Split(lines, " ", -1)
	rndr := &tabulate.PlainRenderer{}
	rndr.SetOFS(" ")
	fmt.Println(rndr.Render(tbl))
}
