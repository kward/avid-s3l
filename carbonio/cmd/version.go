package cmd

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "print carbonio version",
	Long:  `Version prints the build information for the carbonio executable.`,
	Run:   version,
}

func version(cmd *cobra.Command, args []string) {
	fmt.Printf("carbonio version HEAD %s/%s\n", runtime.GOOS, runtime.GOARCH)
}
