package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/Antvirf/metal-control-plane/internal/activities"
	"github.com/Antvirf/metal-control-plane/internal/data"
	"github.com/Antvirf/metal-control-plane/internal/pixieapi"
	"github.com/Antvirf/metal-control-plane/internal/workflows"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

// ONBOARD
// Create workflow
// Take one input, a MAC address
// ARP the MAC to get an IP
// (Ask UNIFI to set this DHCP lease to infinity / create a static IP assignment)
// Query that IP for redhat API info
// Persist to DB: NIC MACs, all info

// OS INSTALL
// Host pixiecore as part of this app
// When a query comes in, based on the MAC of the HOST, construct PXE response --> send to machine for boot!

const (
	QUEUE_NAME = "CONTROL_PLANE_TEMPORAL_QUEUE"
)

func main() {
	mode := flag.String("mode", "onboard", "mode for temporal DC control plane")
	flag.Parse()

	switch *mode {
	case "onboard", "ob":
		c, err := client.Dial(client.Options{})
		if err != nil {
			log.Fatalln("Unable to create Temporal client:", err)
		}
		defer c.Close()

		// Get mac address from first unnamed arg
		flag.Parse()
		if flag.NArg() != 1 {
			log.Fatalf("Incorrect number of positional arguments provided, only one (target MAC) is expected.")
		}
		target_mac_str := flag.Arg(0)

		target_mac, err := net.ParseMAC(target_mac_str)
		if err != nil {
			log.Fatalf("Failed to parse input ('%s') as MAC: %s", target_mac_str, err)
		}
		log.Printf("start onboarding mac: '%s'", target_mac_str)

		request := data.OnboardRequest{
			MacAddress: target_mac,
		}
		options := client.StartWorkflowOptions{
			ID:        fmt.Sprintf("onboard-%s", target_mac),
			TaskQueue: QUEUE_NAME,
		}

		// Send workflow
		we, err := c.ExecuteWorkflow(context.Background(), options, workflows.OnboardMac, request)
		if err != nil {
			log.Fatalln("Unable to start the Workflow:", err)
		}

		log.Printf("WorkflowID: %s RunID: %s\n", we.GetID(), we.GetRunID())

		var result string
		err = we.Get(context.Background(), &result)
		if err != nil {
			log.Fatalln("Unable to get Workflow result:", err)
		}

		log.Println(result)

	case "control-plane", "cp":
		log.Println("start CP server")
		log.Fatal(pixieapi.PixieApiServer().ListenAndServe())

	case "worker":
		c, err := client.Dial(client.Options{})
		if err != nil {
			log.Fatalln("Unable to create Temporal client.", err)
		}
		defer c.Close()
		log.Println("start temporal worker")

		w := worker.New(c, QUEUE_NAME, worker.Options{})
		w.RegisterWorkflow(workflows.OnboardMac)
		w.RegisterActivity(activities.MacToIp)
		w.RegisterActivity(activities.ScrapeFromRedFish)

		err = w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatalln("unable to start Worker", err)
		}
	}
}
