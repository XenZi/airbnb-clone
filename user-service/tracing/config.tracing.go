package tracing

import "os"

type Config struct {
	JaegerAddress string
}

func GetConfig() Config {
	return Config{
		JaegerAddress: os.Getenv("JAEGER_ADDRESS"),
	}
}
