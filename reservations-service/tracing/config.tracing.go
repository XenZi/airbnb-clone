package tracing

import "os"

type Config struct {
	Address                    string
	JaegerAddress              string
	NotificationServiceAddress string
}

func GetConfig() Config {
	return Config{
		Address:                    os.Getenv("RESERVATIONS_SERVICE_ADDRESS"),
		JaegerAddress:              os.Getenv("JAEGER_ADDRESS"),
		NotificationServiceAddress: os.Getenv("NOTIFICATION_SERVICE_ADDRESS"),
	}
}
