package user_reserved

import "metrics-command/commands"

type Reserved struct {
	UserID          string
	AccommodationID string
	ReservedAt      string
}

func NewCommand(userID, accommodationID, reservedAt string) commands.Command {
	return &Reserved{
		UserID:          userID,
		AccommodationID: accommodationID,
		ReservedAt:      reservedAt,
	}
}
