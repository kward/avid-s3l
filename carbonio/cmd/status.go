package cmd

import (
	"fmt"

	"github.com/kward/tabulate/tabulate"
	"github.com/spf13/cobra"
)

var (
	statusCmd = &cobra.Command{
		Use:   "status",
		Short: "device status",
		Long:  `Provide status overview of the carbonio device.`,
		Run:   status,
	}

	statusRaw bool
)

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.Flags().BoolVarP(&statusRaw, "raw", "r", false, "raw output")
}

func status(cmd *cobra.Command, args []string) {
	lines := []string{"LED STATUS"}

	if device == nil {
		fmt.Println("device is uninitialized")
		return
	}

	lines = append(lines, fmt.Sprintf("Power %s", device.LEDs().Power()))
	lines = append(lines, fmt.Sprintf("Status %s", device.LEDs().Status()))
	lines = append(lines, fmt.Sprintf("Mute %s", device.LEDs().Mute()))

	tbl, err := tabulate.NewTable()
	if err != nil {
		fmt.Printf("unable to determine status; %s", err)
		return
	}
	tbl.Split(lines, ifs, -1)
	rndr := &tabulate.PlainRenderer{}
	rndr.SetOFS(ofs)
	fmt.Printf("%s", rndr.Render(tbl))
}
