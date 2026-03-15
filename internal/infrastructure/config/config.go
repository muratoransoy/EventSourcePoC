package config

import "os"

const defaultConnectionString = "kurrentdb://admin:changeit@localhost:2113?tls=false"

type Config struct {
	ConnectionString string
}

func Load() Config {
	connectionString := os.Getenv("KURRENTDB_CONNECTION_STRING")
	if connectionString == "" {
		connectionString = defaultConnectionString
	}

	return Config{
		ConnectionString: connectionString,
	}
}
