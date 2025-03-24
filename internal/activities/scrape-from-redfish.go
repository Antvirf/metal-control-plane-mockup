package activities

import (
	"fmt"
	"net"

	"github.com/Antvirf/metal-control-plane/internal/clients"
	"github.com/Antvirf/metal-control-plane/internal/data"
)

func ScrapeFromRedFish(ip net.IP, mac net.HardwareAddr) (data.HardwareInfo, error) {

	redFishAddress := fmt.Sprintf("http://%s:5000", ip.String())

	info, err := clients.GetRedFishInfo(redFishAddress)
	if err != nil {
		return data.HardwareInfo{}, err
	}

	return data.HardwareInfo{
		BmcIpAddress:  ip,
		BmcMacAddress: mac,
		RedFishData:   info,
	}, nil
}
