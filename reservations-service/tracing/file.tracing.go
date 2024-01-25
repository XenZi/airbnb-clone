package tracing

import "os"

type Config struct {
	Address       string
	JaegerAddress string
}

func GetConfig() Config {
	return Config{
		Address:       os.Getenv("RESERVATIONS_SERVICE_PORT"),
		JaegerAddress: os.Getenv("JAEGER_ADDRESS"),
	}
}
