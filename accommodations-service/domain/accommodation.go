package domain

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"io"
)

type Accommodation struct {
	Id               gocql.UUID
	UserId           string
	Location         string
	Conveniences     string
	MinNumOfVisitors int
	MaxNumOfVisitors int
}

type AccommodationById []*Accommodation

func (aco *Accommodation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(aco)
}

func (aco *Accommodation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(aco)
}

func (ac *AccommodationById) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ac)
}

func (ac *AccommodationById) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ac)
}
