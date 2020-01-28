// The carbionio command enables control of Carbon I/O hardware.
//
// Build for ARM with
// $ GOOS=linux GOARM=7 GOARCH=arm go build carbonio.go
package main

import (
	"fmt"
	"net"

	"github.com/kward/avid-s3l/carbonio/leds"
	"github.com/kward/golib/errors"
	"google.golang.org/grpc/codes"
)

var (
//ip = flag.String("ip")
)

func main() {
	fmt.Println("Hello, world!")
	ip, err := linkLocalIP()
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	fmt.Printf("ip: %v\n", ip)

	for _, led := range []leds.LED{leds.Power, leds.Status, leds.Mute} {
		fmt.Println(led)
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

			if len(ip) != 4 {
				continue
			}
			if ip.IsLinkLocalUnicast() {
				return ip, nil
			}
		}
	}
	return nil, errors.Errorf(codes.NotFound, "unable to determine link local IP address")
}
