package clients

import (
	"context"
	"fmt"
	"os"

	"github.com/unpoller/unifi/v5"
)

func SetupUnifiClient(ctx context.Context) (*unifi.Unifi, error) {
	user := os.Getenv("UNIFI_USERNAME")
	pass := os.Getenv("UNIFI_PASSWORD")
	url := os.Getenv("UNIFI_BASE_URL")
	c := unifi.Config{
		User: user,
		Pass: pass,
		URL:  url,
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
