package domains

type RateAccommodation struct {
	AccommodationID string `json:"accommodationID"`
	Rate            int    `json:"rate"`
	Guest         	Guest `json:"guest"`
}

type RateHost struct {
	Host  Host `json:"host"`
	Rate    int64    `json:"rate"`
	Guest Guest `json:"guest"`
}

type Guest struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
}

type Host struct {
	ID string `json:"id"`
	Email string `json:"email"`
	Username string `json:"username"`
}

type Accommodation struct {
	ID string `json:"id"`
}