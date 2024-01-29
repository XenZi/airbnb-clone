package nats

import (
	"fmt"
	"github.com/nats-io/nats.go"
)

func getConnection(host, port, user, password string) (*nats.Conn, error) {
	url := fmt.Sprintf("nats://%s:%s@%s:%s", user, password, host, port)
	connection, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return connection, nil
}
