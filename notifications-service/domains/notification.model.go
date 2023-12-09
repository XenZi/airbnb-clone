package domains

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserNotification struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserID string `json:"userID" bson:"userID"`
	Notifications []Notification `json:"notifications" bson:"notifications"`
}

type Notification struct {
	Text string `json:"text"`
	CreatedAt string `json:"createdAt"`
	IsOpened bool `json:"isOpened"`
}