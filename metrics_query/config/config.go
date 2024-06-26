package config

import "os"

type Config struct {
	Port      string
	ESDBUser  string
	ESDBPass  string
	ESDBHost  string
	ESDBPort  string
	ESDBGroup string
}

func NewConfig() Config {
	return Config{
		Port:      os.Getenv("PORT"),
		ESDBUser:  os.Getenv("ESDB_USER"),
		ESDBPass:  os.Getenv("ESDB_PASS"),
		ESDBHost:  os.Getenv("ESDB_HOST"),
		ESDBPort:  os.Getenv("ESDB_PORT"),
		ESDBGroup: os.Getenv("ESDB_GROUP"),
	}
}
