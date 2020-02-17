package cmd

import (
	"github.com/kward/avid-s3l/carbonio/servers"
	"github.com/spf13/cobra"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "start the carbonio HTTP and OSC servers",
		Long:  `Server starts the carbonio HTTP and OSC servers.`,
		Run:   server,
	}

	httpPort int
	oscPort  int
)

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVarP(&httpPort, "http_port", "H", 8080, "http port")
	serverCmd.Flags().IntVarP(&oscPort, "osc_port", "O", 41789, "osc port")
}

func server(cmd *cobra.Command, args []string) {
	servers.HttpServer(httpPort, device)
}
