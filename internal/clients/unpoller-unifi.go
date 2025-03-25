package clients

import (
	"context"
	"fmt"

	"github.com/unpoller/unifi/v5"
)

func SetupUnifiClient(ctx context.Context) (*unifi.Unifi, error) {
	c := unifi.Config{
		User: UNIFI_USERNAME,
		Pass: UNIFI_PASSWORD,
		URL:  UNIFI_BASE_URL,
	}

	client, err := unifi.NewUnifi(&c)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate to UniFi: %s", err)
	}
	return client, nil
}

func FindClientDeviceByMac(ctx context.Context, client *unifi.Unifi, mac string) (*unifi.Client, error) {
	sites, err := client.GetSites()
	if err != nil {
		return &unifi.Client{}, fmt.Errorf("failed to authenticate to UniFi: %s", err)
	}

	clientDevices, err := client.GetClients(sites)
	if err != nil {
		return &unifi.Client{}, fmt.Errorf("failed to authenticate to UniFi: %s", err)
	}
	for _, client := range clientDevices {
		if client.Mac == mac {
			return client, nil
		}
	}
	return &unifi.Client{}, fmt.Errorf("No device found with mac %s", mac)

}
