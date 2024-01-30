package domains

type UserJoined struct {
	UserID          string `json:"userID"`
	AccommodationID string `json:"accommodationID"`
	JoinedAt        string `json:"joinedAt"`
	CustomUUID      string `json:"customUUID"`
}

type UserLeft struct {
	UserID          string `json:"userID"`
	AccommodationID string `json:"accommodationID"`
	LeftAt          string `json:"leftAt"`
	CustomUUID      string `json:"customUUID"`
}
