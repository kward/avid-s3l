/*
Package devices enables control of specific Avid S3L devices.
*/
package devices

import (
	"net"

	"github.com/kward/avid-s3l/carbonio/leds"
	"github.com/kward/avid-s3l/carbonio/signals"
	"github.com/kward/golib/errors"
	"google.golang.org/grpc/codes"
)

type Device interface {
	// LED returns defined LEDs.
	LEDs() *leds.LEDs
	// NumMicInputs returns the number of microphone inputs for the device.
	NumMicInputs() int
	// MicInput returns the signal struct for the request input number.
	MicInput(input int) (*signals.Signal, error)
	// IP address of the device.
	IP() net.IP
}

func setDeviceOptions(opts *options) error {
	ip, err := LinkLocalIP()
	if err != nil {
		return errors.Errorf(codes.Internal, "error determining local link ip; %s", err)
	}
	opts.setIP(ip)
	return nil
}

// LinkLocalIP returns the link local IP of the device.
func LinkLocalIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, errors.Errorf(codes.Internal, "unable to enumerate the network interfaces")
	}

	var ip net.IP
	for _, i := range ifaces {
		addrs, err := i.Addrs()
		if err != nil {
			return nil, errors.Errorf(codes.Internal, "unable to enumerate the network addresses")
		}

		for _, addr := range addrs {
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip.IsLinkLocalUnicast() {
				goto found
			}
		}
	}
	ip = nil
found:
	if ip == nil {
		return nil, errors.Errorf(codes.NotFound, "unable to determine link local IP address")
	}
	if ip.Equal(net.ParseIP("fe80::1")) {
		ip = net.ParseIP("127.0.0.1")
	}
	return ip, nil
}
