package activities

import (
	"context"

	"github.com/Antvirf/metal-control-plane/internal/clients"
	"github.com/Antvirf/metal-control-plane/internal/data"
)

func SaveToDb(ctx context.Context, info data.HardwareInfo) (string, error) {
	result, err := clients.WriteHardwareInfo(info)
	if err != nil {
		return "", err
	}
	return result, nil
}
