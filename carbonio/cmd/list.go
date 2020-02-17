package cmd

import (
	"fmt"

	"github.com/kward/avid-s3l/carbonio/signals"
	"github.com/kward/tabulate/tabulate"
	"github.com/spf13/cobra"
)

var (
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "list the carbonio device settings",
		Long:  `List prints the current settings of the carbonio device.`,
		Run:   list,
	}

	listRaw bool
)

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolVarP(&listRaw, "raw", "r", false, "raw output")
}

func list(cmd *cobra.Command, args []string) {
	lines := []string{"SIGNAL GAIN PAD PHANTOM"}

	if device == nil {
		fmt.Println("device is uninitialized")
		return
	}

	for i := uint(1); i <= device.NumMicInputs(); i++ {
		in, err := device.MicInput(i)
		if err != nil {
			fmt.Printf("error accessing mic input %d; %s\n", i, err)
			continue
		}

		gainStr := gain(in)
		padStr := pad(in)
		phantomStr := phantom(in)

		lines = append(lines, fmt.Sprintf("input/mic/%d %s %s %s", i, gainStr, padStr, phantomStr))
	}

	tbl, err := tabulate.NewTable()
	if err != nil {
		fmt.Printf("unable to list settings; %s", err)
		return
	}
	tbl.Split(lines, ifs, -1)
	rndr := &tabulate.PlainRenderer{}
	rndr.SetOFS(ofs)
	fmt.Printf("%s", rndr.Render(tbl))
}

func gain(input *signals.Signal) string {
	if !listRaw {
		v, err := input.Gain()
		if err != nil {
			return "err"
		}
		return fmt.Sprintf("%d", v) // TODO(2020-02-17): include dB unit.
	}
	v, err := input.GainRaw()
	if err != nil {
		return "err"
	}
	return v
}

var boolStr = map[bool]string{
	true:  "On",
	false: "Off",
}

func pad(input *signals.Signal) string {
	if !listRaw {
		v, err := input.Pad()
		if err != nil {
			return "err"
		}
		return boolStr[v]
	}
	v, err := input.PadRaw()
	if err != nil {
		return "err"
	}
	return v
}

func phantom(input *signals.Signal) string {
	if !listRaw {
		v, err := input.Phantom()
		if err != nil {
			return "err"
		}
		return boolStr[v]
	}
	v, err := input.PhantomRaw()
	if err != nil {
		return "err"
	}
	return v
}
