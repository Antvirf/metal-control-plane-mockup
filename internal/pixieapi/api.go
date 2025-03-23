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
)

// TODO: Implement proper server configs lookup, based on data from Redfish
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
		DEFAULT: ServerConfigResponse{
			IpxeScript: "#!ipxe\nchain --autofree http://boot.netboot.xyz/ipxe/netboot.xyz.lkrn",
		},
	}
)

func findConfig(mac string) (ServerType, ServerConfigResponse, error) {
	// Dumb logic against a hardcoded map for now
	serverType, found := SERVER_MAC_TO_TYPE[strings.ToUpper(mac)]
	if !found {
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
	macAddress := filepath.Base(request.URL.Path)

	// Validate that this is a proper MAC, otherwise complain
	_, err := net.ParseMAC(macAddress)
	if err != nil {
		http.Error(writer, "failed to parse MAC address", http.StatusBadRequest)
		return
	}

	// Figure out which boot config to serve
	serverType, config, err := findConfig(macAddress)
	log.Printf("Serving boot config of type '%s' for %s", serverType, filepath.Base(request.URL.Path))
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
