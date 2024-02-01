package user_left

import "metrics-command/commands"

type UserLeftCommand struct {
	UserID          string
	AccommodationID string
	LeftAt          string
	CustomUUID      string
}

func NewCommand(userID, accommodationID, leftAt, customUUID string) commands.Command {
	return &UserLeftCommand{
		UserID:          userID,
		AccommodationID: accommodationID,
		LeftAt:          leftAt,
		CustomUUID:      customUUID,
	}
}
