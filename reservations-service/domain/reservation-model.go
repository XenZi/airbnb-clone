package domain

import (
	"encoding/json"
	"io"

	"github.com/gocql/gocql"
)

type Reservation struct {
	Id                gocql.UUID `json:"id"`
	UserID            string     `json:"userId"`
	AccommodationID   string     `json:"accommodationId"`
	StartDate         string     `json:"startDate"`
	EndDate           string     `json:"endDate"`
	Username          string     `json:"username"`
	AccommodationName string     `json:"accommodationName"`
	Location          string     `json:"location"`
	Price             int        `json:"price"`
	NumberOfDays      int        `json:"numOfDays"`
	Continent         string     `json:"continent"`
}

type FreeReservation struct {
	Id              gocql.UUID `json:"id"`
	AccommodationID string     `json:"accommodationId"`
	StartDate       string     `json:"startDate"`
	EndDate         string     `json:"endDate"`
	Location        string     `json:"location"`
	Price           int        `json:"price"`
	Continent       string     `json:"continent"`
}
type ReservationById []*Reservation

func NewReservation(id gocql.UUID, userID, accommodationID string, startDate, endDate, username, accommodationName, location string, price, numOfDays int, continent string) *Reservation {
	return &Reservation{
		Id:                id,
		UserID:            userID,
		AccommodationID:   accommodationID,
		StartDate:         startDate,
		EndDate:           endDate,
		Username:          username,
		AccommodationName: accommodationName,
		Location:          location,
		Price:             price,
		NumberOfDays:      numOfDays,
		Continent:         continent,
	}

}

func (ac *Reservation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ac)
}

func (ac *Reservation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ac)
}
