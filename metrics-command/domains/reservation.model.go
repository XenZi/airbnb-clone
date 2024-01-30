package domains

type Reservation struct {
	UserID          string `json:"userID"`
	AccommodationID string `json:"accommodationID"`
	ReservedAt      string `json:"reservedAt"`
}
