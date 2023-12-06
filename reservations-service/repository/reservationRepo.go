package repository

import (
	"fmt"
	"log"
	"os"
	"reservation-service/domain"
	"reservation-service/errors"

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
		(id UUID, user_id text, accommodation_id text, start_date text, end_date text, username text, accommodation_name text,location text,price int,num_of_days int,continent text,
		PRIMARY KEY((continent), id)) WITH CLUSTERING ORDER BY (id ASC)`, "reservations")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}

	/*
	   err = rr.session.Query(

	   	fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s

	   (id UUID,accommodation_id text,user_id text, start_date text, end_date,accommodation_name text

	   	PRIMARY KEY((accommodation_id),id))
	   	WITH CLUSTERING ORDER BY(id ASC)`, "reservation_by_accommodation")).Exec()

	   	if err != nil {
	   		rr.logger.Println(err)
	   	}
	*/
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
			(id UUID, accommodation_id string, start_date text, end_date text, location text, price int, continent text,
			 PRIMARY KEY((continent), id))
			WITH CLUSTERING ORDER BY(id ASC)`, "free_accommodation")).Exec()

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

	dropTable("reservations")
	dropTable("free_accommodation")
	// dropTable("reservation_by_accommodation")
	// dropTable("reservation_by_date")
}

func (rr *ReservationRepo) GetReservationsByAccommodation(id string) ([]domain.Reservation, error) {

	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,accommodation_name,location,price,num_of_days FROM reservations WHERE accommodation_id = ? `,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation
		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate, &reservation.EndDate, &reservation.AccommodationName, reservation.Location, reservation.Price, reservation.NumberOfDays)
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

	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,username,accommodation_name,location,price,num_of_days FROM reservations WHERE user_id = ? `,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation
		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate, &reservation.EndDate, &reservation.Username, &reservation.AccommodationName, reservation.Location, reservation.Price, reservation.NumberOfDays)
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

func (rr *ReservationRepo) AvailableDates(accommodationID, startDate, endDate string) (bool, *errors.ReservationError) {
	query := `
    SELECT id, accommodation_id, start_date, end_date, location,price
    FROM free_accommodation 
    WHERE accommodation_id = ? 
    AND start_date <= ? AND end_date >= ?
    ALLOW FILTERING`

	iter := rr.session.Query(query, accommodationID, endDate, startDate).Iter()

	var reservationID string
	if iter.Scan(&reservationID) {
		return false, nil
	}

	if err := iter.Close(); err != nil {
		rr.logger.Println(err)
		return false, errors.NewReservationError(500, "Unable to check availability, database error")
	}
	return true, nil
}

/*
	func (rr *ReservationRepo) InsertReservationByDate(reservation *domain.Reservation) (*domain.Reservation, error) {
		Id, _ := gocql.RandomUUID()
		Quartal, err := utils.GetQuarter(reservation.StartDate)
		if err != nil {
			return nil, err
		}
		err = rr.session.Query(`INSERT INTO reservation_by_date(id,start_date, end_date, quartal, username, accomodation_name) VALUES(?,?,?,?,?,?)`,
			Id, reservation.StartDate, reservation.EndDate, Quartal, reservation.Username, reservation.AccommodationName).Exec()
		if err != nil {
			return nil, err
		}
		reservation.Id = Id
		rr.logger.Println(reservation)
		return reservation, nil
	}
*/
func (rr *ReservationRepo) InsertReservation(reservation *domain.Reservation) (*domain.Reservation, error) {
	Id, _ := gocql.RandomUUID()
	err := rr.session.Query(`INSERT INTO reservations (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,continent)
	VALUES(?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, reservation.StartDate, reservation.EndDate, reservation.Username, reservation.AccommodationName, reservation.Location, reservation.Price, reservation.NumberOfDays, reservation.Continent).Exec()
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}
	reservation.Id = Id
	rr.logger.Println(reservation)
	return reservation, nil

}

/*
	func (rr *ReservationRepo) InsertReservationByAccommodantion(reservation *domain.Reservation) (*domain.Reservation, error) {
		Id, _ := gocql.RandomUUID()
		AccommodationId := "soba1"
		err := rr.session.Query(`INSERT INTO reservation_by_accommodation(id,accommodation_id,user_id,start_date,end_date,accommodation_name)
		VALUES(?,?,?,?,?)`, Id, AccommodationId, reservation.UserID, reservation.StartDate, reservation.EndDate, reservation.AccommodationName).Exec()
		if err != nil {
			rr.logger.Fatalln(err)
			return nil, err
		}
		reservation.Id = Id
		reservation.AccommodationID = AccommodationId
		rr.logger.Println(reservation)
		return reservation, nil
	}
*/
func (rr *ReservationRepo) DeleteById(id string) (*domain.ReservationById, *errors.ReservationError) {
	err := rr.session.Query(`DELETE FROM reservations WHERE id = ?`, id).Exec()
	if err != nil {
		rr.logger.Println(err)
		return nil, errors.NewReservationError(500, "Unable to delete, database error")
	}
	return nil, nil
}
