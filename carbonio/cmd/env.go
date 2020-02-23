package cmd

import (
	"fmt"

	"github.com/kward/avid-s3l/carbonio/devices"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "env",
		Short: "print carbonio environment information",
		Long:  `env prints the environment information of the carbonio executable.`,
		Run:   env,
	})
}

func env(cmd *cobra.Command, args []string) {
	ip, err := devices.LinkLocalIP()
	if err != nil {
		fmt.Printf("link local ip: error %s\n", err)
	} else {
		fmt.Printf("link local ip: %v\n", ip)
	}
}
