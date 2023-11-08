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

//func (ar *AccommodationRepo) CreateTables() {
//	err := ar.session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
//					(id UUID, user_id text, location text, conveniences text, minNumOfVisitors int, maxNumOfVisitors int,
//					PRIMARY KEY ((id), location))
//					WITH CLUSTERING ORDER BY (location ASC)`,
//		"accommodations_by_id")).Exec()
//	if err != nil {
//		ar.logger.Println(err)
//	}
//	err = ar.session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
//					(id UUID, user_id text, location text, conveniences text, minNumOfVisitors int, maxNumOfVisitors int,
//					PRIMARY KEY ((user_id), location))
//					WITH CLUSTERING ORDER BY (location ASC)`,
//		"accommodations_by_user")).Exec()
//	if err != nil {
//		ar.logger.Println(err)
//	}
//
//}

func (ar *AccommodationRepo) CreateTables() {
	err := ar.session.Query(fmt.Sprintf("DROP TABLE IF EXISTS %s", "accommodations_by_id")).Exec()
	if err != nil {
		ar.logger.Println(err)
	}

	err = ar.session.Query(fmt.Sprintf("DROP TABLE IF EXISTS %s", "accommodations_by_user")).Exec()
	if err != nil {
		ar.logger.Println(err)
	}

	err = ar.session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(id UUID, user_id text, location text, conveniences text, minNumOfVisitors int, maxNumOfVisitors int, 
					PRIMARY KEY ((id), location)) 
					WITH CLUSTERING ORDER BY (location ASC)`,
		"accommodations_by_id")).Exec()
	if err != nil {
		ar.logger.Println(err)
	}
	err = ar.session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s 
					(id UUID, user_id text, location text, conveniences text, minNumOfVisitors int, maxNumOfVisitors int, 
					PRIMARY KEY ((user_id), location)) 
					WITH CLUSTERING ORDER BY (location ASC)`,
		"accommodations_by_user")).Exec()
	if err != nil {
		ar.logger.Println(err)
	}
}

func (ar *AccommodationRepo) InsertAccommodationById(accommodation *do.Accommodation) (*do.Accommodation, error) {
	Id, _ := gocql.RandomUUID()
	userId := "Kreirani id"
	ar.logger.Println("Prije kvjerija")
	err := ar.session.Query(`INSERT INTO accommodations_by_id (id,user_id, location, conveniences, minNumOfVisitors, maxNumOfVisitors) 
		VALUES (?, ?, ?, ?, ?, ?)`, Id, userId, accommodation.Location, accommodation.Conveniences, accommodation.MinNumOfVisitors, accommodation.MaxNumOfVisitors).Exec()
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	ar.logger.Println("Poslije kvjerija")
	accommodation.Id = Id
	accommodation.UserId = userId
	ar.logger.Println("Treci kvjerija")
	ar.logger.Println(accommodation)
	return accommodation, nil

}
func (ar *AccommodationRepo) GetAllAccommodations() (do.AccommodationById, error) {
	scanner := ar.session.Query(`SELECT id, user_id, location, conveniences, minNumOfVisitors, maxNumOfVisitors FROM accommodations_by_id`).Iter().Scanner()
	var accommodations do.AccommodationById
	for scanner.Next() {
		var accomm do.Accommodation
		err := scanner.Scan(&accomm.Id, &accomm.UserId, &accomm.Location, &accomm.Conveniences, &accomm.MinNumOfVisitors, &accomm.MaxNumOfVisitors)
		if err != nil {
			ar.logger.Println(err)
			return nil, err
		}
		accommodations = append(accommodations, &accomm)

	}
	if err := scanner.Err(); err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	return accommodations, nil
}

func (ar *AccommodationRepo) GetAccommodationById(id string) (do.AccommodationById, error) {
	scanner := ar.session.Query(`SELECT id, user_id, location, conveniences, minNumOfVisitors, maxNumOfVisitors FROM accommodations_by_id WHERE id=? `, id).Iter().Scanner()
	var accommodations do.AccommodationById
	for scanner.Next() {
		var accomm do.Accommodation
		err := scanner.Scan(&accomm.Id, &accomm.UserId, &accomm.Location, &accomm.Conveniences, &accomm.MinNumOfVisitors, &accomm.MaxNumOfVisitors)
		if err != nil {
			ar.logger.Println(err)
			return nil, err
		}
		accommodations = append(accommodations, &accomm)

	}
	if err := scanner.Err(); err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	return accommodations, nil
}

func (ar *AccommodationRepo) UpdateAccommodationById(id string, location string, accommodation *do.Accommodation) (*do.Accommodation, error) {
	err := ar.session.Query(`UPDATE accommodations_by_id 
                         SET                             
                             conveniences = ?, 
                             minNumOfVisitors = ?, 
                             maxNumOfVisitors = ? 
                         WHERE id = ? AND location=?`,

		accommodation.Conveniences,
		accommodation.MinNumOfVisitors,
		accommodation.MaxNumOfVisitors,
		id, location).Exec()
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	return accommodation, nil

}

func (ar *AccommodationRepo) DeleteAccommodationById(id string) (do.AccommodationById, error) {
	err := ar.session.Query(`DELETE FROM accommodations_by_id WHERE id = ?`, id).Exec()
	if err != nil {
		ar.logger.Println(err)
		return nil, err
	}
	return nil, nil
}
