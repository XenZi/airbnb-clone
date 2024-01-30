package domains

type UserRate struct {
	UserID          string `json:"userID"`
	AccommodationID string `json:"accommodationID"`
	RatedAt         string `json:"ratedAt"`
}
