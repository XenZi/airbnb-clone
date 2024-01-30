package user_joined

import "metrics-command/commands"

type UserJoinedCommand struct {
	UserID          string
	AccommodationID string
	JoinedAt        string
}

func NewCommand(userID, accommodationID, joinedAt string) commands.Command {
	return &UserJoinedCommand{
		UserID:          userID,
		AccommodationID: accommodationID,
		JoinedAt:        joinedAt,
	}
}
