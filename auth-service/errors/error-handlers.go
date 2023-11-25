package errors

import (
	"auth-service/domains"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

func HandleInsertError(err error, user domains.User) (error, int) {
	if writeErr, ok := err.(mongo.WriteException); ok {
		for _, writeError := range writeErr.WriteErrors {
			if writeError.Code == 11000 {
				if strings.Contains(writeError.Message, "email_1") {
					return fmt.Errorf("Duplicate entity with email %s already exists", user.Email), 422
				} else if strings.Contains(writeError.Message, "username_1") {
					return fmt.Errorf("Duplicate entity with username %s already exists", user.Username), 422
				}
			}
		}
	}
	return err, -1
}
