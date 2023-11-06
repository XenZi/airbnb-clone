package repository

import (
	do "accommodations-service/domain"
	"fmt"
	"log"
	"os"

	// NoSQL: module containing Cassandra api client
	"github.com/gocql/gocql"
)

type AccommodationRepo struct {
	session *gocql.Session
	logger  *log.Logger
}

func New(logger *log.Logger) (*AccommodationRepo, error) {
	db := os.Getenv("CASS_DB")

	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	err = session.Query(
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
					WITH replication = {
						'class' : 'SimpleStrategy',
						'replication_factor' : %d
					}`, "student", 1)).Exec()
	if err != nil {
		logger.Println(err)
	}
	session.Close()

	cluster.Keyspace = "student"
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	return &AccommodationRepo{
		session: session,
		logger:  logger,
	}, nil

}

func (ar *AccommodationRepo) CloseSession() {
	ar.session.Close()
}

func (ar *AccommodationRepo) CreateTables() {
	err := ar.session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(id UUID, user_id UUID, location text, conveniences text, minNumOfVisitors int, maxNumOfVisitors int, 
					PRIMARY KEY ((id), location)) 
					WITH CLUSTERING ORDER BY (location ASC)`,
		"accommodations_by_id")).Exec()
	if err != nil {
		ar.logger.Println(err)
	}
	err = ar.session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(id UUID, user_id UUID, location text, conveniences text, minNumOfVisitors int, maxNumOfVisitors int, 
					PRIMARY KEY ((user_id), location)) 
					WITH CLUSTERING ORDER BY (location ASC)`,
		"accommodations_by_user")).Exec()
	if err != nil {
		ar.logger.Println(err)
	}

}

func (ar *AccommodationRepo) InsertAccommodationById(accommodation *do.Accommodation) error {
	Id, _ := gocql.RandomUUID()
	userId, _ := gocql.RandomUUID()
	err := ar.session.Query(`INSERT INTO accommodations_by_id (id,user_id, location, conveniences, minNumOfVisitors, maxNumOfVisitors) 
		VALUES (?, ?, ?, ?, ?, ?)`, Id, userId, accommodation.Location, accommodation.Conveniences, accommodation.MinNumOfVisitors, accommodation.MaxNumOfVisitors).Exec()
	if err != nil {
		ar.logger.Println(err)
		return err
	}
	return nil

}
