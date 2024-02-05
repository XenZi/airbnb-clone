package domains

type RateAccommodation struct {
	AccommodationID string  `json:"accommodationID"`
	Rate            int64   `json:"rate"`
	Guest           Guest   `json:"guest"`
	CreatedAt       string  `json:"createdAt"`
	AvgRating       float64 `json:"avgRating"`
	HostEmail       string  `json:"hostEmail"`
	HostID          string  `json:"hostID"`
}

type RateHost struct {
	Host      Host    `json:"host"`
	Rate      int64   `json:"rate"`
	Guest     Guest   `json:"guest"`
	CreatedAt string  `json:"createdAt"`
	AvgRating float64 `json:"avgRating"`
	HostEmail string  `json:"hostEmail"`
	HostID    string  `json:"hostID"`
}

type Guest struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Host struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string `json:"username"`
}

type Accommodation struct {
	ID string `json:"id"`
}
