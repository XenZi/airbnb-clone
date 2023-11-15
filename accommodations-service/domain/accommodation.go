package domain

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"io"
)

type Accommodation struct {
	Id               gocql.UUID `json:"id"`
	UserId           string     `json:"userId"`
	UserName         string     `json:"username"`
	Name             string     `json:"name"`
	Location         string     `json:"location"`
	Conveniences     string     `json:"conveniences"`
	MinNumOfVisitors int        `json:"minNumOfVisitors"`
	MaxNumOfVisitors int        `json:"maxNumOfVisitors"`
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
