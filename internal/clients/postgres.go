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
		Mac:  input.BmcMacAddress.String(),
		Info: input,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "No match found: %v\n", err)
		os.Exit(1)
	}

	log.Printf("Wrote hardware info to DB for BmcMac '%s'", input.BmcMacAddress.String())
	return record.Mac, nil
}

func GetHardwareInfo(mac net.HardwareAddr) (*data.HardwareInfo, bool, error) {
	conn, err := pgx.Connect(context.Background(), DATABASE_CONNECTION_STRING)
	if err != nil {
		return &data.HardwareInfo{}, false, fmt.Errorf("unable to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	query := sql.New(conn)

	record, err := query.GetHardwareInfo(context.Background(), mac.String())
	if err != nil {
		return &data.HardwareInfo{}, false, nil
	}
	log.Printf("DB successfully queried for hardware info for BmcMac '%s'", mac.String())
	return &record.Info, true, nil
}
