package cmd

import (
	"github.com/kward/avid-s3l/carbonio/handlers"
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
	h := handlers.NewStatusHandler(device)
	h.ServeCommand(cmd.OutOrStdout())
}
