package utils

import (
	"reservation-service/errors"
	"strings"

	"github.com/pariz/gountries"
)

func GetCountry(location string) (string, *errors.ReservationError) {
	locationParts := strings.Split(location, ",")
	if len(locationParts) < 3 {
		return "", errors.NewReservationError(500, "Unable to retrive data")
	}

	country := locationParts[2]
	return country, nil
}

func GetContinent(location string) (string, *errors.ReservationError) {
	country, erro := GetCountry(location)
	if erro != nil {
		return "", errors.NewReservationError(500, "Unable to retrive data")
	}

	countryData := gountries.New()
	result, err := countryData.FindCountryByName(country)
	if err != nil {
		return "", errors.NewReservationError(500, "Unable to retrive data")
	}

	continent := result.Continent
	return continent, nil
}
