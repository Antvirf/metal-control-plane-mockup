# 'Metal Control Plane' POC mockup with Temporal

Basic implementation of [Railway](https://railway.app) ["Zero-Touch Bare Metal at Scale](https://blog.railway.com/p/data-center-build-part-two) article using Go with Temporal for workflow orchestration, Pixiecore for booting, and UniFi for basic network queries.

## Sequence Diagram

The 'bare metal' server in this case can be two separate machines (e.g., VMs) for the purposes of simulating a real server with a BMC.

```mermaid
sequenceDiagram

actor Admin
participant Control Plane
participant Temporal Core
participant Temporal Worker
participant Router
participant Server BMC
participant Server

alt "Machine onboarding flow"
Admin->>Temporal Core: Initiate machine 'onboarding' request of MAC 1 through CLI
Temporal Core->>Temporal Worker: Start 'onboarding workflow
Temporal Worker->>Router: Activity: Check what IP corresponds to MAC 1
Temporal Worker->>Server BMC: Scrape RedFish REST API provided by BMC on given IP
Temporal Worker->>Control Plane: Store server data from RedFish to Control Plane DB
end

alt "Machine OS boot and install workflow"
Admin->>Server: Turn on server/initiate PXE boot (manual in this POC, RedFish controlled in reality)
Server->>Control Plane: Make DHCP boot requests, which are served by Pixiecore
Control Plane->>Control Plane: Check database to determine which boot config to provide
Control Plane->>Server: Pixiecore: Netboot machine with given configs
Server->>Server: Install OS from given configurations
end
```

## Deploy locally

- Make sure your local machine firewall allows UDP traffic on 67,68,69 as these are used by Pixiecore
- Make sure your local router is not serving PXE boot requests, to avoid interfering with the boot process

## Workflow #1 - Onboarding

1. Run RedFish emulator on your 'BMC' host:

    ```bash
    # in case you want to build the container, in case your host is e.g. arm64 like a Raspberry Pi
    podman build -t dmtf/redfish-interface-emulator:latest https://github.com/DMTF/Redfish-Interface-Emulator.git

    podman run --rm --detach --name redfish-emulator -p 5000:5000 dmtf/redfish-interface-emulator:latest
    ```

1. Get the MAC address of the BMC host
1. Configure `.env` with credentials to UniFi router
1. `make containers` - spins up Temporal backend and pixiecore
1. `make worker` - run a Temporal worker
1. Initiate onboarding with `go run main.go onboard $MAC`

This 'onboards' the new machine to the control plane database, with a given hardware config and MAC addresses.

## Workflow #2 - Provisioning/netbooting

1. `make containers` - make sure all containers are up
1. `make control-plane` - run the control plane API
1. Trigger boot - power on a VM for example

The VM will talk to Pixiecore, which in turn asks the control plane which configs to give to the newly booting machine. This config will be provided based on the knowledge the control plane has in its database, given what the onboarding workflow scraped from RedFish.
