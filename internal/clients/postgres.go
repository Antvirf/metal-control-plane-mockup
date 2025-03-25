package clients

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/Antvirf/metal-control-plane/internal/data"
	"github.com/Antvirf/metal-control-plane/internal/sql"
	"github.com/jackc/pgx/v5"
)

func WriteHardwareInfo(input data.HardwareInfo) (string, error) {
	conn, err := pgx.Connect(context.Background(), DATABASE_CONNECTION_STRING)
	if err != nil {
		return "", fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	query := sql.New(conn)

	record, err := query.CreateHardwareInfo(context.Background(), sql.CreateHardwareInfoParams{
		Bmcmac: input.BmcMacAddress.String(),
		Info:   input,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "No match found: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Wrote hardware info to DB for BmcMac '%s'", input.BmcMacAddress.String())
	return record.Bmcmac, nil
}

func GetHardwareInfoByServerMac(mac net.HardwareAddr) (*data.HardwareInfo, bool, error) {
	conn, err := pgx.Connect(context.Background(), DATABASE_CONNECTION_STRING)
	if err != nil {
		return &data.HardwareInfo{}, false, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	query := sql.New(conn)

	// This look up is with the MAC of the *booting* machine, which will not be
	// the same MAC as that of the BMC.
	record, err := query.GetHardwareInfoByEthernetInterfaceMacAddresses(context.Background(), mac.String())
	if err != nil {
		return &data.HardwareInfo{}, false, nil
	}
	log.Printf("DB successfully queried for hardware info for BmcMac '%s'", mac.String())
	return &record.Info, true, nil
}
