package user_reserved

import "metrics-command/commands"

type UserReservedCommand struct {
	UserID          string
	AccommodationID string
	ReservedAt      string
}

func NewCommand(userID, accommodationID, reservedAt string) commands.Command {
	return &UserReservedCommand{
		UserID:          userID,
		AccommodationID: accommodationID,
		ReservedAt:      reservedAt,
	}
}
