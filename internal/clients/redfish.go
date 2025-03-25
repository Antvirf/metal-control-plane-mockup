package clients

import (
	"github.com/Antvirf/metal-control-plane/internal/data"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish"
)

func GetRedFishInfo(address string) (*data.RedFishInfo, error) {
	c, err := gofish.ConnectDefault(address)
	if err != nil {
		panic(err)
	}
	service := c.Service

	// This is a toy implementation, real BMCs may return multiple chassis and systems.
	// All of them should be processed by a real implementation.
	chassis, _ := service.Chassis()
	systems, _ := service.Systems()
	system := systems[0]
	bios, _ := system.Bios()
	processors, _ := system.Processors()
	memory, _ := system.Memory()
	ethernetDevices, _ := system.EthernetInterfaces()
	StorageDevices, _ := system.SimpleStorages()


	// Filter processors to active sockets only
	var activeProcessors []*redfish.Processor
	for _, proc := range processors {
		if proc.Status.State == "Enabled" {
			activeProcessors = append(activeProcessors, proc)
		}
	}


	info := &data.RedFishInfo{
		Chassis:            chassis[0],
		System:             system,
		Bios:               bios,
		Processors:         activeProcessors,
		Memory:             memory,
		EthernetInterfaces: ethernetDevices,
		StorageDevices:     StorageDevices,
	}

	return info, nil
}
