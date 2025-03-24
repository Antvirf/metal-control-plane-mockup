package clients

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/Antvirf/metal-control-plane/internal/data"
	"github.com/Antvirf/metal-control-plane/internal/sql"
	"github.com/jackc/pgx/v5"
)

func WriteHardwareInfo(input data.HardwareInfo) (string, error) {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://controlplane:controlplane@localhost/controlplane"
	}

	conn, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
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

	fmt.Println("wrote to DB!")
	fmt.Println(record.Mac)
	fmt.Println(record.Info.BmcIpAddress)
	fmt.Println(record.Info.RedFishData.Bios)
	return record.Mac, nil
}

func GetHardwareInfo(mac net.HardwareAddr) {
	dbUrl := os.Getenv("DATABASE_URL")
	if dbUrl == "" {
		dbUrl = "postgres://controlplane:controlplane@localhost/controlplane"
	}

	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	query := sql.New(conn)

	record, err := query.GetHardwareInfo(context.Background(), mac.String())
	if err != nil {
		fmt.Fprintf(os.Stderr, "No match found: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(record.Mac)
	fmt.Println(record.Info.BmcIpAddress)
	fmt.Println(record.Info.RedFishData.Bios)
}
