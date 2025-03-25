package pixieapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Antvirf/metal-control-plane/internal/clients"
)

//go:generate stringer -type=Pill
type ServerType int

const (
	ST_COMPUTE_G1 ServerType = iota
	ST_COMPUTE_G2
	DEFAULT
	IGNORE
)

func findConfig(mac net.HardwareAddr) (ServerType, ServerConfigResponse, error) {
	// Get machine info from DB, and then figure out what boot configs
	// the machine should get based on its hardware configuration.
	// Note that this lookup is by the MAC address of the *booting* machine, which
	// will be one of the NICs of the server - *NOT* the BMC MAC address.
	hardwareInfo, found, err := clients.GetHardwareInfoByServerMac(mac)
	var serverType ServerType
	if err != nil {
		// In case of an actual error, return "IGNORE"
		return IGNORE, ServerConfigResponse{}, fmt.Errorf("error querying db: %v", err)
	}
	if !found {
		// If no error, but nothing was found, what do you want to do here? Serve DEFAULT, or IGNORE?
		// For my case, IGNORE. I only want to PXE things that exist in my DB.
		return IGNORE, ServerConfigResponse{}, errors.New("DB query returned no results, this machine has not been onboarded")
	}

	// Choosing the server type based on the returned value, toy example
	processorModel := strings.ToLower(hardwareInfo.RedFishData.Processors[0].Model)
	switch {
	case strings.Contains(processorModel, "intel"):
		serverType = ST_COMPUTE_G1
	case strings.Contains(processorModel, "amd"):
		serverType = ST_COMPUTE_G2
	default:
		serverType = DEFAULT
	}

	config, found := SERVER_TYPE_TO_CONFIG[serverType]
	if !found {
		return DEFAULT, ServerConfigResponse{}, fmt.Errorf("server type found, but no config defined for type %s", serverType)
	}
	return serverType, config, nil
}

func pixieApiHandler(writer http.ResponseWriter, request *http.Request) {
	// Extract the last part of the path - this is a mac address
	macAddressRaw := filepath.Base(request.URL.Path)

	// Validate that this is a proper MAC, otherwise complain
	macAddress, err := net.ParseMAC(macAddressRaw)
	if err != nil {
		http.Error(writer, "failed to parse MAC address", http.StatusBadRequest)
		return
	}

	// Figure out which boot config to serve
	serverType, config, err := findConfig(macAddress)
	if err != nil {
		log.Printf("[%s] error looking up config for server: %v", macAddressRaw, err)
	}
	log.Printf("[%s] serving boot config of type '%s'", macAddressRaw, serverType)
	if err := json.NewEncoder(writer).Encode(&config); err != nil {
		panic(err)
	}
}

func PixieApiServer() *http.Server {
	http.HandleFunc("/v1/boot/", pixieApiHandler)

	return &http.Server{
		Addr:           ":8081",
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
}
