package cmd

import (
	"fmt"

	"github.com/kward/avid-s3l/carbonio/handlers"
	"github.com/kward/avid-s3l/carbonio/helpers"
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
	h, err := handlers.NewHandlers(device,
		handlers.Raw(listRaw))
	if err != nil {
		helpers.Exit(fmt.Sprintf("error instantiating handlers; %s", err))
	}
	h.ListCommand(cmd.OutOrStdout())
}
