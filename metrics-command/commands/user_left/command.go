package user_left

import "metrics-command/commands"

type UserLeftCommand struct {
	UserID          string
	AccommodationID string
	LeftAt          string
}

func NewCommand(userID, accommodationID, leftAt string) commands.Command {
	return &UserLeftCommand{
		UserID:          userID,
		AccommodationID: accommodationID,
		LeftAt:          leftAt,
	}
}
