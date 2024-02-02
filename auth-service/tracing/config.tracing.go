package tracing

import "os"

type Config struct {
	Address            string
	JaegerAddress      string
	UserServiceAddress string
}

func GetConfig() Config {
	return Config{
		Address:            os.Getenv("AUTH_SERVICE_ADDRESS"),
		JaegerAddress:      os.Getenv("JAEGER_ADDRESS"),
		UserServiceAddress: os.Getenv("USER_SERVICE_ADDRESS"),
	}
}
