package pixieapi

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/Antvirf/metal-control-plane/internal/clients"
)

type ServerConfigResponse struct {
	Kernel     string   `json:"kernel,omitempty"`
	Initrd     []string `json:"initrd,omitempty"`
	Cmdline    string   `json:"cmdline,omitempty"`
	IpxeScript string   `json:"ipxe-script,omitempty"`
	// https://github.com/danderson/netboot/blob/main/pixiecore/README.api.md#custom-ipxe-boot-script
}

//go:generate stringer -type=Pill
type ServerType int

const (
	ST_COMPUTE_G1 ServerType = iota
	ST_COMPUTE_G2
	DEFAULT
	IGNORE
)

// Hardcoded map of SERVER TYPE --> Boot configs
// You'd probably want to control this from a DB in a real implementation, or from a mountable config file.
var (
	SERVER_MAC_TO_TYPE = map[string]ServerType{
		// Must provide MACs in upper case.
		"BC:24:11:17:55:95": ST_COMPUTE_G1,
	}
	SERVER_TYPE_TO_CONFIG = map[ServerType]ServerConfigResponse{
		ST_COMPUTE_G1: ServerConfigResponse{
			Kernel: "https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os/images/pxeboot/vmlinuz",
			Initrd: []string{
				"https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os/images/pxeboot/initrd.img",
			},
			Cmdline: "selinux=0 inst.repo=https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os inst.text",
		},
		ST_COMPUTE_G2: ServerConfigResponse{
			Kernel: "https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os/images/pxeboot/vmlinuz",
			Initrd: []string{
				"https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os/images/pxeboot/initrd.img",
			},
			Cmdline: "selinux=1 inst.repo=https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os inst.text",
		},
		DEFAULT: ServerConfigResponse{
			IpxeScript: "#!ipxe\nchain --autofree http://boot.netboot.xyz/ipxe/netboot.xyz.lkrn",
		},
	}
)

func findConfig(mac net.HardwareAddr) (ServerType, ServerConfigResponse, error) {
	var serverType ServerType
	// Implementation #1: Hardcoded logic against a map
	// serverType, found := SERVER_MAC_TO_TYPE[strings.ToUpper(mac)]
	// if !found {
	// 	serverType = DEFAULT
	// }

	// Implementation #2: Get machine info from DB, and then figure out what boot configs
	// the machine should get based on its hardware configuration.
	// Note that this lookup is by the MAC address of the *booting* machine, which
	// will be one of the NICs of the server - *NOT* the BMC MAC address.
	hardwareInfo, found, err := clients.GetHardwareInfoByServerMac(mac)
	if err != nil {
		// In case of an actual error, return "IGNORE"
		return IGNORE, ServerConfigResponse{}, fmt.Errorf("error querying db: %v", err)
	}
	if !found {
		// If no error, but nothing was found, what do you want to do here? Serve DEFAULT, or IGNORE?
		// For my case, IGNORE. I only want to PXE things that exist in my DB.
		return IGNORE, ServerConfigResponse{}, fmt.Errorf("DB query returned no results, this machine has not been onboarded", err)
	}

	// Choosing the server type based on the returned value, toy example
	processorModel := strings.ToLower(hardwareInfo.RedFishData.Processors[0].Model)
	switch {
	case strings.Contains(processorModel, "intel"):
		log.Printf("assigning mac=%s with type %s", mac.String(), ST_COMPUTE_G1)
		serverType = ST_COMPUTE_G1
	case strings.Contains(processorModel, "amd"):
		log.Printf("assigning mac=%s with type %s", mac.String(), ST_COMPUTE_G2)
		serverType = ST_COMPUTE_G2
	default:
		log.Printf("assigning mac=%s with type %s", mac.String(), DEFAULT)
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
	return
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
