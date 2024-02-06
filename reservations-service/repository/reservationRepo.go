package repository

import (
	"context"
	"fmt"
	"log"

	"os"
	"reservation-service/config"
	"reservation-service/domain"
	"reservation-service/errors"
	"reservation-service/utils"

	"time"

	"github.com/gocql/gocql"
	"github.com/pariz/gountries"
	"go.opentelemetry.io/otel/trace"
)

type ReservationRepo struct {
	session *gocql.Session
	logger  *config.Logger
	tracer  trace.Tracer
}

// db config and creating keyspace
func New(logger *config.Logger, tracer trace.Tracer) (*ReservationRepo, error) {
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
		tracer:  tracer,
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
			(id UUID, accommodation_id text,  location text, price int, continent text, country text, date_range set<text>,
			 PRIMARY KEY((accommodation_id),country,id))
			WITH CLUSTERING ORDER BY(country ASC,id ASC)`, "free_accommodation")).Exec()

	if err != nil {
		rr.logger.Println(err)
	}
	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
			(id UUID, accommodation_id text,  location text, price int, continent text, country text, date_range set<text>,is_active boolean,
			 PRIMARY KEY((is_active),price,id))
			WITH CLUSTERING ORDER BY(price ASC,id ASC)`, "avl_by_price")).Exec()

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
	) WITH CLUSTERING ORDER BY ( id ASC)`, "reservation_by_user")).Exec()

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
		PRIMARY KEY (host_id,user_id,end_date, id)
	) WITH CLUSTERING ORDER BY (user_id ASC,end_date ASC, id ASC)`, "reservation_by_host")).Exec()

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
			PRIMARY KEY (accommodation_id,user_id,end_date, id)
		) WITH CLUSTERING ORDER BY (user_id ASC,end_date ASC, id ASC)`, "reservation_by_accommodation")).Exec()

	if err != nil {
		rr.logger.Println(err)
	}

	err = rr.session.Query(
		fmt.Sprintf(`CREATE TABLE IF NOT EXISTS %s
			(id UUID, host_id text, ,
			 PRIMARY KEY((host_id),id))
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
	//dropTable("free_accommodation")
	dropTable("reservation_by_user")
	dropTable("reservation_by_host")
	dropTable("reservation_by_accommodation")
	dropTable("deleted_reservations")
	//	dropTable("avl_by_price")

}

func (rr *ReservationRepo) GetReservationsByUser(ctx context.Context, id string) ([]domain.Reservation, error) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.GetReservationsByUser")
	defer span.End()
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
			rr.logger.LogError("reservationsRepo", err.Error())
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	if err := scanner.Err(); err != nil {
		rr.logger.LogError("reservationsRepo", err.Error())
		return nil, err
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Found reservations by userID: %v", reservations))
	return reservations, nil
}

func (rr *ReservationRepo) GetReservationsByHost(ctx context.Context, id string) ([]domain.Reservation, error) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.GetReservationsByHost")
	defer span.End()
	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,username,accommodation_name,location,price,
	num_of_days,date_range,is_active,country,host_id FROM reservation_by_host
	 WHERE  host_id = ?`,
		id).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation

		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate,
			&reservation.EndDate, &reservation.Username, &reservation.AccommodationName, &reservation.Location, &reservation.Price,
			&reservation.NumberOfDays, &reservation.DateRange, &reservation.IsActive, &reservation.Country, &reservation.HostID)
		if err != nil {
			rr.logger.LogError("reservationsRepo", err.Error())
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	if err := scanner.Err(); err != nil {
		rr.logger.LogError("reservationsRepo", err.Error())
		return nil, err
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Found reservations by hostID: %v", reservations))
	return reservations, nil
}

func (rr *ReservationRepo) InsertAvailability(ctx context.Context, reservation *domain.FreeReservation) (*domain.FreeReservation, error) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.InsertAvailability")
	defer span.End()
	country, err := utils.GetCountry(reservation.Location)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	continent, err := utils.GetContinent(reservation.Location)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	for _, drwp := range reservation.DateRange {
		batch := rr.session.NewBatch(gocql.LoggedBatch)
		ID, _ := gocql.RandomUUID()
		batch.Query(`
				INSERT INTO free_accommodation (id, accommodation_id, location, price, continent, country, date_range)
				VALUES(?, ?, ?, ?, ?, ?, ?)
			`, ID, reservation.AccommodationID, reservation.Location, drwp.Price, continent, country, drwp.DateRange)
		batch.Query(`
				INSERT INTO avl_by_price (id, accommodation_id, location, price, continent, country, date_range,is_active)
				VALUES(?, ?, ?, ?, ?, ?, ?,?)
			`, ID, reservation.AccommodationID, reservation.Location, drwp.Price, continent, country, drwp.DateRange, true)

		if err := rr.session.ExecuteBatch(batch); err != nil {
			rr.logger.LogError("reservationsRepo", err.Error())
			return nil, err
		}

	}

	reservation.Continent = continent
	reservation.Country = country
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Inserted availability: %v", reservation))
	return reservation, nil
}

func (rr *ReservationRepo) AvailableDates(ctx context.Context, accommodationID string, dateRange []string) ([]domain.FreeReservation, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.AvailableDates")
	defer span.End()
	var result []domain.FreeReservation
	for _, date := range dateRange {
		query := `
        SELECT id, accommodation_id,location, price,country
        FROM free_accommodation 
        WHERE accommodation_id = ? 
        AND date_range CONTAINS ?
        `

		iter := rr.session.Query(query, accommodationID, date).Iter()

		var reservation domain.FreeReservation
		for iter.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.Location, &reservation.Price, reservation.Country) {
			result = append(result, reservation)
		}

		if err := iter.Close(); err != nil {
			rr.logger.LogError("reservationsRepo", err.Error())
			return nil, errors.NewReservationError(500, "Unable to check availability, database error")
		}
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Found availability by accommodationID and dateRange: %v", result))
	return result, nil
}

func (rr *ReservationRepo) InsertReservation(ctx context.Context, reservation *domain.Reservation) (*domain.Reservation, error) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.InsertReservation")
	defer span.End()
	Id, _ := gocql.RandomUUID()
	country, err := utils.GetCountry(reservation.Location)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}

	continent, err := utils.GetContinent(reservation.Location)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	log.Println("COUNTRY< COUNTINET", country, continent)
	startDate := reservation.DateRange[0]
	endDate := reservation.DateRange[len(reservation.DateRange)-1]
	println(startDate, endDate)

	batch := rr.session.NewBatch(gocql.LoggedBatch)

	// Insert into reservations table
	batch.Query(`INSERT INTO reservations (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,
	    continent,date_range,is_active,country,host_id)
	    VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, startDate,
		endDate, reservation.Username, reservation.AccommodationName, reservation.Location,
		reservation.Price, reservation.NumberOfDays, continent, reservation.DateRange, true, country, reservation.HostID)
	batch.Query(`INSERT INTO reservation_by_user (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,
	    continent,date_range,is_active,country,host_id)
	    VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, startDate,
		endDate, reservation.Username, reservation.AccommodationName, reservation.Location,
		reservation.Price, reservation.NumberOfDays, continent, reservation.DateRange, true, country, reservation.HostID)
	batch.Query(`INSERT INTO reservation_by_host (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,
	    continent,date_range,is_active,country,host_id)
	    VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, startDate,
		endDate, reservation.Username, reservation.AccommodationName, reservation.Location,
		reservation.Price, reservation.NumberOfDays, continent, reservation.DateRange, true, country, reservation.HostID)
	batch.Query(`INSERT INTO reservation_by_accommodation (id,user_id,accommodation_id,start_date,end_date,username,accommodation_name,location,price,num_of_days,
			continent,date_range,is_active,country,host_id)
			VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)`, Id, reservation.UserID, reservation.AccommodationID, startDate,
		endDate, reservation.Username, reservation.AccommodationName, reservation.Location,
		reservation.Price, reservation.NumberOfDays, continent, reservation.DateRange, true, country, reservation.HostID)

	if err := rr.session.ExecuteBatch(batch); err != nil {
		rr.logger.LogError("reservationsRepo", err.Error())
		return nil, err
	}

	reservation.Id = Id
	reservation.Country = country
	reservation.Continent = continent
	reservation.IsActive = true
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Inserted reservation: %v", reservation))

	return reservation, nil
}

func (rr *ReservationRepo) DeleteById(ctx context.Context, country string, id, userID, hostID, accommodationID, endDate string) (*domain.Reservation, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.DeleteById")
	defer span.End()
	countryData := gountries.New()

	result, err := countryData.FindCountryByName(country)
	if err != nil {
		return nil, errors.NewReservationError(500, err.Error())
	}
	continent := result.Continent
	batch := rr.session.NewBatch(gocql.LoggedBatch)

	batch.Query(`DELETE FROM reservations WHERE continent = ? AND country = ? AND id = ?`, continent, country, id)
	batch.Query(`DELETE FROM reservation_by_user WHERE user_id = ? AND id = ?`, userID, id)
	batch.Query(`DELETE FROM reservation_by_host WHERE host_id = ? AND user_id = ? AND end_date = ? AND id = ? `, hostID, userID, endDate, id)
	batch.Query(`DELETE FROM reservation_by_accommodation WHERE accommodation_id = ? AND user_id = ? AND end_date = ? AND id = ?`, accommodationID, userID, endDate, id)
	batch.Query(`INSERT INTO deleted_reservations(id,host_id) VALUES(?,?)`, id, hostID)

	if err := rr.session.ExecuteBatch(batch); err != nil {
		rr.logger.LogError("reservationsRepo", err.Error())
		return nil, errors.NewReservationError(500, "Unable to cancel the reservation")
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Deleted reservation by ID: %v", id))

	return nil, nil
}

func (rr *ReservationRepo) ReservationsInDateRange(ctx context.Context, accommodationIDs []string, dateRange []string) ([]string, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.ReservationsInDateRange")
	defer span.End()
	uniqueReservations := make(map[string]struct{})
	for _, date := range dateRange {
		for _, accommodationID := range accommodationIDs {
			query := `SELECT accommodation_id FROM reservation_by_accommodation WHERE accommodation_id = ? AND date_range CONTAINS ?`

			iter := rr.session.Query(query, accommodationID, date).Iter()

			var reservation string
			for iter.Scan(&reservation) {
				uniqueReservations[reservation] = struct{}{}
			}
			if err := iter.Close(); err != nil {
				rr.logger.LogError("reservationsRepo", err.Error())
				return nil, errors.NewReservationError(500, "Unable to retrieve reservations, database error")
			}
		}
	}

	result := make([]string, 0, len(uniqueReservations))
	for key := range uniqueReservations {
		result = append(result, key)
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Found reservation by accommodationIDs and dateRange: %v", result))

	return result, nil
}

func (rr *ReservationRepo) IsAvailable(ctx context.Context, accommodationID string, dateRange []string) (bool, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.IsAvailable")
	defer span.End()
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
			rr.logger.LogError("reservationsRepo", err.Error())
			return false, errors.NewReservationError(500, "Unable to check availability, database error")
		}
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Found out is accommodation available or not by accommodationID and dateRange: %v", false))

	return false, nil
}

func (rr *ReservationRepo) CheckAvailabilityForAccommodation(ctx context.Context, accommodationID string) ([]domain.GetAvailabilityForAccommodation, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.CheckAvailabilityForAccommodation")
	defer span.End()
	var result []domain.GetAvailabilityForAccommodation

	query := `
    SELECT date_range,price,id
    FROM free_accommodation 
    WHERE accommodation_id = ? 
    `

	iter := rr.session.Query(query, accommodationID).Iter()
	var dateRange []string
	var price int
	var id string
	var avl domain.GetAvailabilityForAccommodation
	for iter.Scan(&dateRange, &price, &id) {
		avl.DateRange = dateRange
		avl.Price = price
		avl.Id = id
		result = append(result, avl)
	}
	if err := iter.Close(); err != nil {
		rr.logger.LogError("reservationsRepo", err.Error())
		return nil, errors.NewReservationError(500, "Unable to check availability, database error")
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Checked availability by accommodationID : %v", result))

	return result, nil

}

func (rr *ReservationRepo) IsReserved(ctx context.Context, accommodationID string, dateRange []string) (bool, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.IsReserved")
	defer span.End()

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
			rr.logger.LogError("reservationsRepo", err.Error())
			return false, errors.NewReservationError(500, "Unable to check is reserved, database error")
		}

	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Checked if is accommodation reserved by accommodationID and dateRanges: %v", false))

	return false, nil
}
func (rr *ReservationRepo) GetNumberOfCanceledReservations(ctx context.Context, hostID string) (int, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.GetNumberOfCanceledReservations")
	defer span.End()
	query := `SELECT COUNT(*) FROM deleted_reservations WHERE host_id = ?`
	iter := rr.session.Query(query, hostID).Iter()

	var numberOfCanceled int
	if iter.Scan(&numberOfCanceled) {
		return numberOfCanceled, nil
	}
	return 0, errors.NewReservationError(500, "Failed to get the number of canceled reservations")
}

func (rr *ReservationRepo) GetTotalReservationsByHost(ctx context.Context, hostID string) (int, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.GetTotalReservationsByHost")
	defer span.End()
	query := `SELECT COUNT(*) FROM reservation_by_host WHERE host_id = ?`
	iter := rr.session.Query(query, hostID).Iter()
	var totalReservations int
	if iter.Scan(&totalReservations) {
		return totalReservations, nil
	}
	return 0, errors.NewReservationError(500, "Failed to get the number of total reservations")

}
func (rr *ReservationRepo) GetReservationsByAccommodationWithEndDate(ctx context.Context, accommodationID, userID string) ([]domain.Reservation, error) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.GetReservationsByAccommodationWithEndDate")
	defer span.End()
	currentDate := time.Now().Format("2006-01-02")
	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,username,accommodation_name,location,price,
	num_of_days,date_range,is_active,country,host_id FROM reservation_by_accommodation
	 WHERE  accommodation_id = ? AND user_id = ? AND end_date <= ?`,
		accommodationID, userID, currentDate).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation

		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate,
			&reservation.EndDate, &reservation.Username, &reservation.AccommodationName, &reservation.Location, &reservation.Price,
			&reservation.NumberOfDays, &reservation.DateRange, &reservation.IsActive, &reservation.Country, &reservation.HostID)
		if err != nil {
			rr.logger.LogError("reservationsRepo", err.Error())
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	if err := scanner.Err(); err != nil {
		rr.logger.LogError("reservationsRepo", err.Error())
		return nil, err
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Found expired reservations by accommodationID: %v", reservations))
	return reservations, nil

}

func (rr *ReservationRepo) GetReservationsByHostWithEndDate(ctx context.Context, hostID, userID string) ([]domain.Reservation, error) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.GetReservationsByHostWithEndDate")
	defer span.End()
	currentDate := time.Now().Format("2006-01-02")
	scanner := rr.session.Query(`SELECT id,accommodation_id, user_id, start_date, end_date,username,accommodation_name,location,price,
	num_of_days,date_range,is_active,country,host_id FROM reservation_by_host
	 WHERE  host_id = ? AND user_id = ? AND end_date <= ?`,
		hostID, userID, currentDate).Iter().Scanner()

	var reservations []domain.Reservation
	for scanner.Next() {
		var reservation domain.Reservation

		err := scanner.Scan(&reservation.Id, &reservation.AccommodationID, &reservation.UserID, &reservation.StartDate,
			&reservation.EndDate, &reservation.Username, &reservation.AccommodationName, &reservation.Location, &reservation.Price,
			&reservation.NumberOfDays, &reservation.DateRange, &reservation.IsActive, &reservation.Country, &reservation.HostID)
		if err != nil {
			rr.logger.LogError("reservationsRepo", err.Error())
			return nil, err
		}

		reservations = append(reservations, reservation)
	}

	if err := scanner.Err(); err != nil {
		rr.logger.LogError("reservationsRepo", err.Error())
		return nil, err
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Found expired reservations by hostID: %v", reservations))
	return reservations, nil

}

func (rr *ReservationRepo) DeleteAvl(ctx context.Context, accommodationID, id, country string, price int) (*domain.FreeReservation, error) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.DeleteAvl")
	defer span.End()

	batch := rr.session.NewBatch(gocql.LoggedBatch)

	batch.Query(`DELETE FROM free_accommodation WHERE accommodation_id = ? AND country = ? AND id = ?`, accommodationID, country, id)
	batch.Query(`DELETE FROM avl_by_price WHERE is_active = ? AND price = ? AND id = ?`, true, price, id)

	if err := rr.session.ExecuteBatch(batch); err != nil {
		rr.logger.LogError("reservationsRepo", err.Error())
		return nil, errors.NewReservationError(500, "Unable to delete availability")
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Deleted availability by ID: %v", id))
	return nil, nil

}

func (rr *ReservationRepo) GetAccommodationIDsByMaxPrice(ctx context.Context, maxPrice int) ([]string, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.GetAccommodationIDsByMaxPrice")
	defer span.End()
	scanner := rr.session.Query(`
        SELECT accommodation_id FROM avl_by_price WHERE is_active = ? AND price <= ?
    `, true, maxPrice).Iter().Scanner()

	var accommodationIDs []string
	for scanner.Next() {
		var accommodationID string
		err := scanner.Scan(&accommodationID)
		if err != nil {
			rr.logger.LogError("reservationsRepo", err.Error())
			return nil, errors.NewReservationError(500, "Unable to retrive the data")
		}
		accommodationIDs = append(accommodationIDs, accommodationID)

	}

	if erro := scanner.Err(); erro != nil {
		rr.logger.LogError("reservationsRepo", erro.Error())
		return nil, errors.NewReservationError(500, "Unable to retrive the data")
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Found accommodations by price: %v", accommodationIDs))
	return accommodationIDs, nil
}

func (rr *ReservationRepo) AvailabilityNotInDateRange(ctx context.Context, accommodationIDs []string, dateRange []string) ([]string, *errors.ReservationError) {
	ctx, span := rr.tracer.Start(ctx, "ReservationRepo.AvailabilityNotInDateRange")
	defer span.End()
	uniqueAccommodationIDs := make(map[string]struct{})

	for _, accommodationID := range accommodationIDs {
		query := `SELECT accommodation_id FROM free_accommodation WHERE accommodation_id = ? AND date_range CONTAINS ?`

		isInDateRange := false

		for _, date := range dateRange {
			iter := rr.session.Query(query, accommodationID, date).Iter()

			var result string
			if iter.Scan(&result) {
				isInDateRange = true
			}

			if err := iter.Close(); err != nil {
				rr.logger.LogError("reservationsRepo", err.Error())
				return nil, errors.NewReservationError(500, "Unable to retrieve availability, database error")
			}
		}
		if !isInDateRange {
			uniqueAccommodationIDs[accommodationID] = struct{}{}
		}
	}

	result := make([]string, 0, len(uniqueAccommodationIDs))
	for key := range uniqueAccommodationIDs {
		result = append(result, key)
	}
	rr.logger.LogInfo("reservationRepo", fmt.Sprintf("Found availabilities that are not in date range by accommodationIDs and dateRange: %v", result))
	return result, nil
}
