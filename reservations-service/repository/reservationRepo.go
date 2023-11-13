package repository

import (
	"fmt"
	"log"
	"os"
	"reservation-service/domain"

	"github.com/gocql/gocql"
)

type ReservationRepo struct {
	session *gocql.Session
	logger  *log.Logger
}

// db config and creating keyspace
func New(logger *log.Logger) (*ReservationRepo, error) {
	db := os.Getenv("CASS_DB")

	cluster := gocql.NewCluster(db)
	cluster.Keyspace = "system"
	session, err := cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	err = session.Query(fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s
	WITH replication = {
		'class' : 'SimpleStrategy',
		'replication_factor' : %d
	}`, "reservation", 1)).Exec()

	if err != nil {
		logger.Println(err)
	}
	session.Close()

	cluster.Keyspace = "reservation"
	cluster.Consistency = gocql.One
	session, err = cluster.CreateSession()
	if err != nil {
		logger.Println(err)
		return nil, err
	}

	return &ReservationRepo{
		session: session,
		logger:  logger,
	}, nil
}

func (rr *ReservationRepo) CloseSession() {
	rr.session.Close()
}

func (rr *ReservationRepo) CreateTables() {
	// Drop existing tables first
	rr.DropTables()

	// Create new tables
	err := rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
		(id UUID,user_id text, accommodation_id text, start_date text, end_date text, 
			PRIMARY KEY((user_id),id)) WITH CLUSTERING ORDER BY (accommodation_id DESC,start_date ASC)`, "reservation_by_user")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
	(id UUID,accommodation_id text,user_id text, start_date text, end_date text,
		PRIMARY KEY((accommodation_id),id, user_id, start_date))
		WITH CLUSTERING ORDER BY(user_id DESC,start_date ASC)`, "reservation_by_accommodation")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}
}

func (rr *ReservationRepo) DropTables() {
	dropTable := func(tableName string) {
		err := rr.session.Query(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName)).Exec()
		if err != nil {
			rr.logger.Println(err)
		} else {
			rr.logger.Printf("Table %s dropped successfully", tableName)
		}
	}

	dropTable("reservation_by_user")
	dropTable("reservation_by_accommodation")
}

func (rr *ReservationRepo) GetReservationsByAccommodation(id string) ([]domain.Reservation, error) {

	scanner := rr.session.Query(`SELECT accommodation_id, user_id, start_date, end_date FROM reservation_by_accommodation WHERE accommodation_id = ? `,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation
		err := scanner.Scan(&reservation.Id)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return reservations, nil
}

func (rr *ReservationRepo) GetReservationsByUser(id string) ([]domain.Reservation, error) {

	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date FROM reservation_by_user WHERE user_id = ? `,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation
		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate, &reservation.EndDate)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	return reservations, nil
}

func (rr *ReservationRepo) InsertReservationByUser(reservation *domain.Reservation) (*domain.Reservation, error) {
	Id, _ := gocql.RandomUUID()
	UserId := "myID"
	err := rr.session.Query(`INSERT INTO reservation_by_user(id,user_id,accommodation_id,start_date,end_date)
	VALUES(?,?,?,?,?)`, Id, UserId, reservation.AccommodationID, reservation.StartDate, reservation.EndDate).Exec()
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	reservation.Id = Id
	reservation.UserID = UserId
	rr.logger.Println(reservation)
	return reservation, nil

}

func (rr *ReservationRepo) InsertReservationByAccommodantion(reservation *domain.Reservation) (*domain.Reservation, error) {
	Id, _ := gocql.RandomUUID()
	AccommodationId := "soba1"
	err := rr.session.Query(`INSERT INTO reservation_by_accommodation(id,accommodation_id,user_id,start_date,end_date)
	VALUES(?,?,?,?,?)`, Id, AccommodationId, reservation.UserID, reservation.StartDate, reservation.EndDate).Exec()
	if err != nil {
		rr.logger.Fatalln(err)
		return nil, err
	}
	reservation.Id = Id
	reservation.AccommodationID = AccommodationId
	rr.logger.Println(reservation)
	return reservation, nil
}
