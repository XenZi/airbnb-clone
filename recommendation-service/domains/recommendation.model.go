package domains

type Recommendation struct {
	AccommodationID string  `json:"accommodationID"`
	Rating          float64 `json:"rating"`
}
