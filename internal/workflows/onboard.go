package workflows

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/Antvirf/metal-control-plane/internal/activities"
	"github.com/Antvirf/metal-control-plane/internal/data"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func OnboardMac(ctx workflow.Context, req data.OnboardRequest) (string, error) {
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    3 * time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    10 * time.Second,
		MaximumAttempts:    500, // 0 is unlimited retries
	}
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retrypolicy,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	// Activity #1: Get IP from Unifi
	var targetIp net.IP
	ipMappingError := workflow.ExecuteActivity(ctx, activities.MacToIp, req).Get(ctx, &targetIp)
	if ipMappingError != nil {
		return "", ipMappingError
	}

	// Activity #2: Scrape this IP for information from RedFish
	var info data.HardwareInfo
	bmcScrapingError := workflow.ExecuteActivity(ctx, activities.ScrapeFromRedFish, targetIp, req.MacAddress).Get(ctx, &info)
	if bmcScrapingError != nil {
		return "", bmcScrapingError
	}

	// Activity #3: Persist in DB
	var dbResult string
	dbSaveError := workflow.ExecuteActivity(ctx, activities.SaveToDb, info).Get(ctx, &dbResult)
	if dbSaveError != nil {
		return "", dbSaveError
	}

	result := fmt.Sprintf(
		"Found %s/%s, system ID: %s, processor: %s, memory: %f GiB",
		info.BmcIpAddress,
		info.BmcMacAddress,
		info.RedFishData.System.ID,
		info.RedFishData.Processors[0].Model,
		info.RedFishData.System.MemorySummary.TotalSystemMemoryGiB,
	)

	log.Printf("Machine '%s' onboarded successfully, IP: '%s'", req.MacAddress, targetIp)
	return result, nil
}
