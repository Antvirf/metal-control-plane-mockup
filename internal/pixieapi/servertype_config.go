package pixieapi

// Hardcoded map of SERVER TYPE --> Boot configs
// You'd probably want to control this from a DB in a real implementation, or from a mountable config file.
var (
	SERVER_TYPE_TO_CONFIG = map[ServerType]ServerConfigResponse{
		ST_COMPUTE_G1: {
			Kernel: "https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os/images/pxeboot/vmlinuz",
			Initrd: []string{
				"https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os/images/pxeboot/initrd.img",
			},
			Cmdline: "selinux=0 inst.repo=https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os inst.text",
		},
		ST_COMPUTE_G2: {
			Kernel: "https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os/images/pxeboot/vmlinuz",
			Initrd: []string{
				"https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os/images/pxeboot/initrd.img",
			},
			Cmdline: "selinux=1 inst.repo=https://mirror.stream.centos.org/9-stream/BaseOS/x86_64/os inst.text",
		},
		DEFAULT: {
			IpxeScript: "#!ipxe\nchain --autofree http://boot.netboot.xyz/ipxe/netboot.xyz.lkrn",
		},
	}
)

type ServerConfigResponse struct {
	Kernel     string   `json:"kernel,omitempty"`
	Initrd     []string `json:"initrd,omitempty"`
	Cmdline    string   `json:"cmdline,omitempty"`
	IpxeScript string   `json:"ipxe-script,omitempty"`
	// https://github.com/danderson/netboot/blob/main/pixiecore/README.api.md#custom-ipxe-boot-script
}
