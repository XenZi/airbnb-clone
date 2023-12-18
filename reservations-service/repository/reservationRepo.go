package repository

import (
	"fmt"
	"log"
	"os"
	"reservation-service/domain"
	"reservation-service/errors"
	"strings"

	"github.com/gocql/gocql"
	"github.com/pariz/gountries"
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
		(id UUID, user_id text, accommodation_id text, start_date text, end_date text, username text, accommodation_name text,location text,price int,num_of_days int,continent text,date_range text,is_active boolean,country text,host_id text,
		PRIMARY KEY((continent),country,id)) WITH CLUSTERING ORDER BY (country ASC,id ASC)`, "reservations")).Exec()
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
			(id UUID, accommodation_id text, start_date text, end_date text, location text, price int, continent text,country text,host_id text,
			 PRIMARY KEY((continent),country, id))
			WITH CLUSTERING ORDER BY(country ASC,id ASC)`, "free_accommodation")).Exec()

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

	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,accommodation_name,location,price,
	num_of_days,date_range,country,host_id FROM reservations WHERE is_active = true AND accommodation_id = ?  ALLOW FILTERING`,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation
		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate, &reservation.EndDate,
			&reservation.AccommodationName, &reservation.Location, &reservation.Price, &reservation.NumberOfDays, &reservation.DateRange,
			&reservation.Country, &reservation.HostID)
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
	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,username,accommodation_name,location,
	price,num_of_days,date_range,is_active,country,host_id FROM reservations WHERE is_active = true AND user_id = ? ALLOW FILTERING `,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation
		var dateRangeString string

		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate, &reservation.EndDate, &reservation.Username,
			&reservation.AccommodationName, &reservation.Location, &reservation.Price, &reservation.NumberOfDays,
			&dateRangeString, &reservation.IsActive, &reservation.Country, &reservation.HostID)
		if err != nil {
			rr.logger.Println(err)
			return nil, err
		}

		reservation.DateRange = strings.Split(dateRangeString, ",")

		reservations = append(reservations, reservation)
	}

	if err := scanner.Err(); err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	return reservations, nil
}

func (rr *ReservationRepo) InsertAvailability(reservation *domain.FreeReservation) (*domain.FreeReservation, error) {
	Id, _ := gocql.RandomUUID()
	countryData := gountries.New()
	locationParts := strings.Split(reservation.Location, ",")
	if len(locationParts) < 3 {
		return nil, errors.NewReservationError(400, "Invalid location format: %s")
	}

	country := locationParts[2]
	result, err := countryData.FindCountryByName(country)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	continent := result.Continent

	err = rr.session.Query(`
		INSERT INTO free_accommodation (id, accommodation_id, start_date, end_date, location, price, continent,country,host_id)
		VALUES(?, ?, ?, ?, ?, ?, ?,?,?)
	`, Id, reservation.AccommodationID, reservation.StartDate, reservation.EndDate, reservation.Location, reservation.Price, continent, country, reservation.HostID).Exec()
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	reservation.Id = Id
	reservation.Continent = continent
	reservation.Country = country
	return reservation, nil
}

func (rr *ReservationRepo) AvailableDates(accommodationID, startDate, endDate string) ([]domain.FreeReservation, *errors.ReservationError) {
	var result []domain.FreeReservation

	query := `
        SELECT id, accommodation_id, start_date, end_date, location, price,country,host_id,continent
        FROM free_accommodation 
        WHERE accommodation_id = ? 
        AND start_date <= ? AND end_date >= ?`

	iter := rr.session.Query(query, accommodationID, startDate, endDate).Iter()

	var reservation domain.FreeReservation
	for iter.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.StartDate, &reservation.EndDate, &reservation.Location,
		&reservation.Price, &reservation.Country, &reservation.HostID, &reservation.Continent) {
		result = append(result, reservation)
	}

	if err := iter.Close(); err != nil {
		rr.logger.Println(err)
		return nil, errors.NewReservationError(500, "Unable to check availability, database error")
	}

	return result, nil
}

func (rr *ReservationRepo) InsertReservation(reservation *domain.Reservation) (*domain.Reservation, error) {
	Id, _ := gocql.RandomUUID()
	//StartDate := reservation.DateRange[0]
	//EndDate := reservation.DateRange[len(reservation.DateRange)-1]

	dateRangeString := strings.Join(reservation.DateRange, ",")

	countryData := gountries.New()
	locationParts := strings.Split(reservation.Location, ",")
	if len(locationParts) < 3 {
		return nil, errors.NewReservationError(400, "Invalid location format: %s")
	}

	country := locationParts[2]
	result, err := countryData.FindCountryByName(country)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	continent := result.Continent

	err = rr.session.Query(`INSERT INTO reservations (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,
		continent,date_range,is_active,country,host_ids)
	VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, reservation.StartDate, reservation.EndDate, reservation.Username,
		reservation.AccommodationName, reservation.Location, reservation.Price, reservation.NumberOfDays, continent, dateRangeString, true, country, reservation.HostID).Exec()
	if err != nil {
		rr.logger.Println(err)
		return nil, err
	}

	reservation.Id = Id
	reservation.Country = country
	reservation.Continent = continent
	reservation.IsActive = true
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
func (rr *ReservationRepo) DeleteById(country string, id string) (*domain.Reservation, *errors.ReservationError) {
	countryData := gountries.New()

	result, err := countryData.FindCountryByName(country)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	continent := result.Continent
	err = rr.session.Query(`UPDATE reservations SET is_active = false WHERE continent = ? AND country = ? AND id = ?`, continent, country, id).Exec()
	if err != nil {
		rr.logger.Println(err)
		return nil, errors.NewReservationError(500, "Unable to delete, database error")
	}

	return nil, nil
}
func (rr *ReservationRepo) ReservationsInDateRange(accommodationIDs []string, dateRange []string) ([]string, *errors.ReservationError) {
	var result []string

	query := `
        SELECT id, accommodation_id, continent, country
        FROM reservations 
        WHERE accommodation_id IN ? 
        AND start_date <= ? AND end_date >= ?`

	iter := rr.session.Query(query, accommodationIDs, dateRange[0], dateRange[len(dateRange)-1]).Iter()

	var reservationID, accommodationID, continent, country string
	for iter.Scan(&reservationID, &accommodationID, &continent, &country) {
		result = append(result, reservationID)
	}

	if err := iter.Close(); err != nil {
		rr.logger.Println(err)
		return nil, errors.NewReservationError(500, "Unable to retrieve reservations, database error")
	}

	return result, nil
}

func (rr *ReservationRepo) IsAvailable(accommodationID, startDate, endDate string) (bool, *errors.ReservationError) {
	query := `
	SELECT id, accommodation_id, continent, country
	FROM free_accommodation 
	WHERE accommodation_id = ? 
	AND start_date <= ? AND end_date >= ?`

	iter := rr.session.Query(query, accommodationID, startDate, endDate).Iter()

	var reservationID, continent, country string
	if iter.Scan(&reservationID, &accommodationID, &continent, &country) {
		return false, nil
	}

	if err := iter.Close(); err != nil {
		rr.logger.Println(err)
		return false, errors.NewReservationError(500, "Unable to check availability, database error")
	}
	return true, nil
}

func (rr *ReservationRepo) CheckAvailabilityForAccommodation(accommodationID string) ([]domain.GetAvailabilityForAccommodation, *errors.ReservationError) {
	var result []domain.GetAvailabilityForAccommodation

	query := `
    SELECT start_date, end_date,price,id,continent,country
    FROM free_accommodation 
    WHERE accommodation_id = ? `

	iter := rr.session.Query(query, accommodationID).Iter()
	var startDate string
	var endDate string
	var price int
	var id string
	var continent string
	var country string
	var avl domain.GetAvailabilityForAccommodation
	for iter.Scan(&startDate, &endDate, &price, &id, &continent, &country) {
		avl.StartDate = startDate
		avl.EndDate = endDate
		avl.Price = price
		result = append(result, avl)
	}
	if err := iter.Close(); err != nil {
		rr.logger.Println(err)
		return nil, errors.NewReservationError(500, "Unable to check availability, database error")
	}

	return result, nil

}

func (rr *ReservationRepo) IsReserved(accommodationID, startDate, endDate string) (bool, *errors.ReservationError) {
	query := `
	SELECT id,continent,country
	FROM reservations 
	WHERE accommodation_id = ? 
	AND start_date <= ? AND end_date >= ?`

	iter := rr.session.Query(query, accommodationID, startDate, endDate).Iter()

	var reservationID string
	var continent string
	var country string
	if iter.Scan(&reservationID, &continent, &country) {
		return true, nil
	}

	if err := iter.Close(); err != nil {
		rr.logger.Println(err)
		return false, errors.NewReservationError(500, "Unable to check is reserved, database error")
	}
	return false, nil
}
