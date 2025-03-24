package data

import (
	"net"

	"github.com/stmcginnis/gofish/redfish"
)

type OnboardRequest struct {
	MacAddress net.HardwareAddr
}

type RedFishInfo struct {
	Chassis            *redfish.Chassis
	System             *redfish.ComputerSystem
	Bios               *redfish.Bios
	Processors         []*redfish.Processor
	Memory             []*redfish.Memory
	EthernetInterfaces []*redfish.EthernetInterface
	StorageDevices     []*redfish.SimpleStorage
}

type HardwareInfo struct {
	BmcIpAddress  net.IP
	BmcMacAddress net.HardwareAddr
	RedFishData   *RedFishInfo
}
