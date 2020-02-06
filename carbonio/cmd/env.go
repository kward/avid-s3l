package cmd

import (
	"fmt"
	"net"

	"github.com/kward/golib/errors"
	"github.com/spf13/cobra"
	"google.golang.org/grpc/codes"
)

func init() {
	rootCmd.AddCommand(envCmd)
}

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "print carbonio environment information",
	Long:  `Env prints the environment information of the carbonio executable.`,
	Run:   env,
}

func env(cmd *cobra.Command, args []string) {
	ip, err := linkLocalIP()
	if err != nil {
		fmt.Printf("link local ip: error %s\n", err)
	} else {
		fmt.Printf("link local ip: %v\n", ip)
	}
}

// linkLocalIP returns the link local IP of the device.
func linkLocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.Errorf(codes.Internal, "unable to enumerate the network interfaces")
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, errors.Errorf(codes.Internal, "unable to enumerate the network addresses")
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsLinkLocalUnicast() {
				return ip, nil
			}
		}
	}
	return nil, errors.Errorf(codes.NotFound, "unable to determine link local IP address")
}
