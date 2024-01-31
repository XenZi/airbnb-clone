package tracing

import "os"

type Config struct {
	ServiceName   string
	JaegerAddress string
}

func GetConfig() Config {
	return Config{
		ServiceName:   "reservations-service",
		JaegerAddress: os.Getenv("JAEGER_ADDRESS"),
	}
}
