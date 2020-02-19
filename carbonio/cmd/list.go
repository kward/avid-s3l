package cmd

import (
	"github.com/kward/avid-s3l/carbonio/handlers"
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
	h := handlers.NewListHandler(device, handlers.Raw(listRaw))
	h.ServeCommand(cmd.OutOrStdout())
}
