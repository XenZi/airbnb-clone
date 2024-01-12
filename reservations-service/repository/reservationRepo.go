package repository

import (
	"fmt"
	"log"
	"os"
	"reservation-service/domain"
	"reservation-service/errors"
	"reservation-service/utils"

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
		(id UUID, user_id text, accommodation_id text, start_date text, end_date text, username text, accommodation_name text,location text,price int,
			num_of_days int,continent text, date_range set<text>,is_active boolean,country text,host_id text,
		PRIMARY KEY((continent),country,id)) WITH CLUSTERING ORDER BY (country ASC,id ASC)`, "reservations")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
			(id UUID, accommodation_id text, start_date text, end_date text, location text, price int, continent text, country text, date_range set<text>,
			 PRIMARY KEY((accommodation_id),country,id))
			WITH CLUSTERING ORDER BY(country ASC,id ASC)`, "free_accommodation")).Exec()

	if err != nil {
		rr.logger.Println(err)
	}
	err = rr.session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id UUID,
		user_id text,
		accommodation_id text,
		start_date text,
		end_date text,
		username text,
		accommodation_name text,
		location text,
		price int,
		num_of_days int,
		continent text,
		date_range set<text>,
		is_active boolean,
		country text,
		host_id text,
		PRIMARY KEY (user_id, id)
	) WITH CLUSTERING ORDER BY (id ASC)`, "reservation_by_user")).Exec()

	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
		id UUID,
		user_id text,
		accommodation_id text,
		start_date text,
		end_date text,
		username text,
		accommodation_name text,
		location text,
		price int,
		num_of_days int,
		continent text,
		date_range set<text>,
		is_active boolean,
		country text,
		host_id text,
		PRIMARY KEY (host_id, id)
	) WITH CLUSTERING ORDER BY (id ASC)`, "reservation_by_host")).Exec()

	if err != nil {
		rr.logger.Println(err)
	}
	err = rr.session.Query(fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s (
			id UUID,
			user_id text,
			accommodation_id text,
			start_date text,
			end_date text,
			username text,
			accommodation_name text,
			location text,
			price int,
			num_of_days int,
			continent text,
			date_range set<text>,
			is_active boolean,
			country text,
			host_id text,
			PRIMARY KEY (accommodation_id, id)
		) WITH CLUSTERING ORDER BY ( id ASC)`, "reservation_by_accommodation")).Exec()

	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
			(id UUID, accommodation_id text, ,
			 PRIMARY KEY((accommodation_id),id))
			WITH CLUSTERING ORDER BY(id ASC)`, "deleted_reservations")).Exec()

	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(fmt.Sprintf("CREATE INDEX ON reservation_by_accommodation (date_range);")).Exec()
	if err != nil {
		rr.logger.Println(err)
	}
	err = rr.session.Query(fmt.Sprintf("CREATE INDEX ON free_accommodation (date_range);")).Exec()
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
	dropTable("reservation_by_user")
	dropTable("reservation_by_host")
	dropTable("reservation_by_accommodation")
	dropTable("deleted_reservations")

}

/*func (rr *ReservationRepo) GetReservationsByAccommodation(id string) ([]domain.Reservation, error) {

	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,accommodation_name,location,price,
	num_of_days,date_range,country,continent,host_id,is_active FROM reservations WHERE is_active = true AND accommodation_id = ? `,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation
		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate, &reservation.EndDate,
			&reservation.AccommodationName, &reservation.Location, &reservation.Price, &reservation.NumberOfDays,
			&reservation.DateRange, &reservation.Country, &reservation.Continent, &reservation.HostID, &reservation.IsActive)
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
}*/

func (rr *ReservationRepo) GetReservationsByUser(id string) ([]domain.Reservation, error) {
	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,username,accommodation_name,location,price,
	num_of_days,date_range,is_active,country,host_id FROM reservation_by_user
	 WHERE user_id = ?`,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation

		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate,
			&reservation.EndDate, &reservation.Username, &reservation.AccommodationName, &reservation.Location, &reservation.Price,
			&reservation.NumberOfDays, &reservation.DateRange, &reservation.IsActive, &reservation.Country, &reservation.HostID)
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
func (rr *ReservationRepo) GetReservationsByHost(id string) ([]domain.Reservation, error) {
	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,username,accommodation_name,location,price,
	num_of_days,date_range,is_active,country,host_id FROM reservation_by_host
	 WHERE  host_id = ?`,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation

		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate,
			&reservation.EndDate, &reservation.Username, &reservation.AccommodationName, &reservation.Location, &reservation.Price,
			&reservation.NumberOfDays, reservation.DateRange, &reservation.IsActive, &reservation.Country, &reservation.HostID)
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
func (rr *ReservationRepo) InsertAvailability(reservation *domain.FreeReservation) (*domain.FreeReservation, error) {
	country, err := utils.GetCountry(reservation.Location)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	continent, err := utils.GetContinent(reservation.Location)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	for _, dateRange := range reservation.DateRange {
		Id, _ := gocql.RandomUUID()
		log.Println(dateRange)
		query := rr.session.Query(`
			INSERT INTO free_accommodation (id, accommodation_id, start_date, end_date, location, price, continent, country, date_range)
			VALUES(?, ?, ?, ?, ?, ?, ?, ?, ?)
		`, Id, reservation.AccommodationID, reservation.StartDate, reservation.EndDate, reservation.Location, reservation.Price, continent, country, dateRange)

		if err := query.Exec(); err != nil {
			rr.logger.Println(err)
			return nil, errors.NewReservationError(500, err.Error())
		}
	}
	reservation.Continent = continent
	reservation.Country = country
	return reservation, nil
}

func (rr *ReservationRepo) AvailableDates(accommodationID string, dateRange []string) ([]domain.FreeReservation, *errors.ReservationError) {
	var result []domain.FreeReservation
	for _, date := range dateRange {
		query := `
        SELECT id, accommodation_id, start_date, end_date, location, price,country
        FROM free_accommodation 
        WHERE accommodation_id = ? 
        AND date_range CONTAINS ?
        `

		iter := rr.session.Query(query, accommodationID, date).Iter()

		var reservation domain.FreeReservation
		for iter.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.StartDate, &reservation.EndDate, &reservation.Location, &reservation.Price, reservation.Country) {
			result = append(result, reservation)
		}

		if err := iter.Close(); err != nil {
			rr.logger.Println(err)
			return nil, errors.NewReservationError(500, "Unable to check availability, database error")
		}
	}

	return result, nil
}

func (rr *ReservationRepo) InsertReservation(reservation *domain.Reservation) (*domain.Reservation, error) {
	Id, _ := gocql.RandomUUID()
	// dateRangeString := strings.Join(reservation.DateRange, ",")
	// dateRangeS
	country, err := utils.GetCountry(reservation.Location)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	continent, err := utils.GetContinent(reservation.Location)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	log.Println("COUNTRY< COUNTINET", country, continent)

	batch := rr.session.NewBatch(gocql.LoggedBatch)

	// Insert into reservations table
	batch.Query(`INSERT INTO reservations (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,
	    continent,date_range,is_active,country,host_id)
	    VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, reservation.StartDate,
		reservation.EndDate, reservation.Username, reservation.AccommodationName, reservation.Location,
		reservation.Price, reservation.NumberOfDays, continent, reservation.DateRange, true, country, reservation.HostID)
	batch.Query(`INSERT INTO reservation_by_user (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,
	    continent,date_range,is_active,country,host_id)
	    VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, reservation.StartDate,
		reservation.EndDate, reservation.Username, reservation.AccommodationName, reservation.Location,
		reservation.Price, reservation.NumberOfDays, continent, reservation.DateRange, true, country, reservation.HostID)
	batch.Query(`INSERT INTO reservation_by_host (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,
	    continent,date_range,is_active,country,host_id)
	    VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, reservation.StartDate,
		reservation.EndDate, reservation.Username, reservation.AccommodationName, reservation.Location,
		reservation.Price, reservation.NumberOfDays, continent, reservation.DateRange, true, country, reservation.HostID)
	batch.Query(`INSERT INTO reservation_by_accommodation (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,
			continent,date_range,is_active,country,host_id)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, reservation.StartDate,
		reservation.EndDate, reservation.Username, reservation.AccommodationName, reservation.Location,
		reservation.Price, reservation.NumberOfDays, continent, reservation.DateRange, true, country, reservation.HostID)

	if err := rr.session.ExecuteBatch(batch); err != nil {
		//	rr.logger.Println("ASDASDDSAASD ", err)
		return nil, err
	}

	reservation.Id = Id
	reservation.Country = country
	reservation.Continent = continent
	reservation.IsActive = true
	rr.logger.Println(reservation)

	return reservation, nil
}

func (rr *ReservationRepo) DeleteById(country string, id, userID, hostID, accommodationID string) (*domain.Reservation, *errors.ReservationError) {
	countryData := gountries.New()

	result, err := countryData.FindCountryByName(country)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	continent := result.Continent
	batch := rr.session.NewBatch(gocql.LoggedBatch)

	batch.Query(`DELETE FROM reservations WHERE continent = ? AND country = ? AND id = ?`, continent, country, id)
	batch.Query(`DELETE FROM reservation_by_user WHERE user_id = ? AND id = ?`, userID, id)
	batch.Query(`DELETE FROM reservation_by_host WHERE host_id = ? AND id = ?`, hostID, id)
	batch.Query(`DELETE FROM reservation_by_accommodation WHERE accommodation_id = ? AND id = ?`, accommodationID, id)
	batch.Query(`INSERT INTO deleted_reservations(id,accommodation_id) VALUES(?,?)`, id, accommodationID)

	if err := rr.session.ExecuteBatch(batch); err != nil {
		rr.logger.Println(err)
		return nil, errors.NewReservationError(500, "Unable to cancel the reservation")
	}

	return nil, nil
}

func (rr *ReservationRepo) ReservationsInDateRange(accommodationIDs []string, dateRange []string) ([]string, *errors.ReservationError) {
	var result []string
	for _, date := range dateRange {
		query := `
        SELECT accommodation_id
        FROM reservation_by_accommodation 
        WHERE accommodation_id IN ? 
        AND date_range CONTAINS ? 
    `

		iter := rr.session.Query(query, accommodationIDs, date).Iter()

		var reservation string
		for iter.Scan(&reservation) {
			result = append(result, reservation)
		}

		if err := iter.Close(); err != nil {
			rr.logger.Println(err)
			return nil, errors.NewReservationError(500, "Unable to retrieve reservations, database error")
		}
	}

	return result, nil
}

func (rr *ReservationRepo) IsAvailable(accommodationID string, dateRange []string) (bool, *errors.ReservationError) {
	for _, date := range dateRange {
		query := `
	SELECT id
	FROM free_accommodation 
	WHERE accommodation_id = ? 
	AND date_range CONTAINS ?
	`

		iter := rr.session.Query(query, accommodationID, date).Iter()

		var reservationID string
		if iter.Scan(&reservationID) {
			return true, nil
		}

		if err := iter.Close(); err != nil {
			rr.logger.Println(err)
			return false, errors.NewReservationError(500, "Unable to check availability, database error")
		}
	}

	return false, nil
}

func (rr *ReservationRepo) CheckAvailabilityForAccommodation(accommodationID string) ([]domain.GetAvailabilityForAccommodation, *errors.ReservationError) {
	var result []domain.GetAvailabilityForAccommodation

	query := `
    SELECT start_date, end_date,price
    FROM free_accommodation 
    WHERE accommodation_id = ? 
    `

	iter := rr.session.Query(query, accommodationID).Iter()
	var startDate string
	var endDate string
	var price int
	var avl domain.GetAvailabilityForAccommodation
	for iter.Scan(&startDate, &endDate, &price) {
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

func (rr *ReservationRepo) IsReserved(accommodationID string, dateRange []string) (bool, *errors.ReservationError) {

	for _, date := range dateRange {
		query := `
		SELECT id
		FROM reservation_by_accommodation
		WHERE accommodation_id = ? 
		AND date_range CONTAINS ?
		`
		iter := rr.session.Query(query, accommodationID, date).Iter()

		var reservationID string
		if iter.Scan(&reservationID) {
			return true, nil
		}

		if err := iter.Close(); err != nil {
			rr.logger.Println(err)
			return false, errors.NewReservationError(500, "Unable to check is reserved, database error")
		}

	}
	return false, nil
}
