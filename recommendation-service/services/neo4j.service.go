package services

import (
	"log"
	"os"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

type Neo4j struct {
	driver neo4j.DriverWithContext
}

func NewNeo4j() (*Neo4j, error) {
	uri := os.Getenv("NEO4J_DB")
	user := os.Getenv("NEO4J_USERNAME")
	pass := os.Getenv("NEO4J_PASS")
	auth := neo4j.BasicAuth(user, pass, "")
	driver, err := neo4j.NewDriverWithContext(uri, auth)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &Neo4j{
		driver: driver,
	}, nil
}

func(n Neo4j) GetDriver() neo4j.DriverWithContext {
	return n.driver
}