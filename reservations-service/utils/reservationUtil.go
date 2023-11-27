package utils

import (
	"fmt"
	"time"
)

func GetQuarter(startDate string) (string, error) {
	date, err := time.Parse("2023-11-06", startDate)
	if err != nil {
		return "", fmt.Errorf("error parsing date: %v", err)
	}

	month := date.Month()

	var quarter string
	switch {
	case month >= 1 && month <= 3:
		quarter = "Q1"
	case month >= 4 && month <= 6:
		quarter = "Q2"
	case month >= 7 && month <= 9:
		quarter = "Q3"
	case month >= 10 && month <= 12:
		quarter = "Q4"
	default:
		return "", fmt.Errorf("invalid month")
	}

	return quarter, nil
}
