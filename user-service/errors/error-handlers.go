package errors

import (
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
	"user-service/domain"
)

func HandleInsertError(err error, user domain.User) (error, int) {
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

func HandleNoDocumentsError(err error, id string) (error, int) {
	if err == mongo.ErrNoDocuments {
		return fmt.Errorf("No entity with id %s found", id), 404
	}
	return err, -1
}
