package user_rated

import "metrics-command/commands"

type UserRatedCommand struct {
	UserID                  string
	AccommodationID         string
	RatedAt                 string
	expectedLastEventNumber int64
	number                  uint64
}

func NewCommand(userID, accommodationID, ratedAt string) commands.Command {
	return &UserRatedCommand{
		UserID:          userID,
		AccommodationID: accommodationID,
		RatedAt:         ratedAt,
	}
}
