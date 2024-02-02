package tracing

import "os"

type Config struct {
	Address            string
	JaegerAddress      string
	UserServiceAddress string
}

func GetConfig() Config {
	return Config{
		Address:            os.Getenv("RESERVATIONS_SERVICE_ADDRESS"),
		JaegerAddress:      os.Getenv("JAEGER_ADDRESS"),
		UserServiceAddress: os.Getenv("NOTIFICATION_SERVICE_ADDRESS"),
	}
}
