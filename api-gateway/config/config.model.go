package config

import "os"

type Config struct {
	Address string
	AuthServiceAddress string
}

func GetConfig() Config {
	return Config{
		Address: os.Getenv("GATEWAY_ADDRESS"),
		AuthServiceAddress: os.Getenv("AUTH_SERVICE_ADDRESS"),
	}
}