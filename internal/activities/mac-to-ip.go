package activities

import (
	"context"
	"fmt"
	"net"

	"github.com/Antvirf/metal-control-plane/internal/clients"
	"github.com/Antvirf/metal-control-plane/internal/data"
)

func MacToIp(ctx context.Context, req data.OnboardRequest) (net.IP, error) {
	unifiClient, err := clients.SetupUnifiClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to set up unifi client: %s", err)
	}

	device, err := clients.FindClientDeviceByMac(ctx, unifiClient, req.MacAddress.String())
	if err != nil {
		return nil, fmt.Errorf("failed to get device from UniFi: %s", err)
	}

	addr := net.ParseIP(device.IP)
	if addr == nil { // Failed to parse
		return nil, fmt.Errorf("could not parse device IP: %s", device.IP)
	}
	return addr, nil
}
