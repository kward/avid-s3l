package cmd

import (
	"fmt"

	"github.com/kward/avid-s3l/carbonio/handlers"
	"github.com/kward/avid-s3l/carbonio/helpers"
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
	h, err := handlers.NewHandlers(device,
		handlers.Raw(listRaw))
	if err != nil {
		helpers.Exit(fmt.Sprintf("error instantiating handlers; %s", err))
	}
	h.StatusCommand(cmd.OutOrStdout())
}
