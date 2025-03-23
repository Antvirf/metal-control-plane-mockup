package data

import "net"

type OnboardRequest struct {
	MacAddress net.HardwareAddr
}
