package domain

import (
	"encoding/json"
	"io"

	"github.com/gocql/gocql"
)

type Reservation struct {
	Id              gocql.UUID `json: "id"`
	UserID          string     `json: "userId"`
	AccommodationID string     `json: "accommodationId"`
	StartDate       string     `json: "startDate"`
	EndDate         string     `json: "endDate"`
}

func NewReservation(id gocql.UUID, userID, accommodationID string, startDate, endDate string) *Reservation {
	return &Reservation{
		Id:              id,
		UserID:          userID,
		AccommodationID: accommodationID,
		StartDate:       startDate,
		EndDate:         endDate,
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
