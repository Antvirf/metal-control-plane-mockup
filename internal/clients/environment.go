package clients

import "os"

var (
	DATABASE_CONNECTION_STRING = os.Getenv("DATABASE_URL")
	UNIFI_USERNAME             = os.Getenv("UNIFI_USERNAME")
	UNIFI_PASSWORD             = os.Getenv("UNIFI_PASSWORD")
	UNIFI_BASE_URL             = os.Getenv("UNIFI_BASE_URL")
)
