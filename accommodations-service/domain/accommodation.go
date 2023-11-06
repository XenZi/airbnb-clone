package domain

import (
	"encoding/json"
	"github.com/gocql/gocql"
	"io"
)

type Accommodation struct {
	Id               gocql.UUID
	Owner            User
	Location         string
	Conveniences     string
	MinNumOfVisitors int
	MaxNumOfVisitors int
}

func (ac *Accommodation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ac)
}

func (ac *Accommodation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ac)
}
