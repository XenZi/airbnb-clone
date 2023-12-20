package domains

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserNotification struct {
	ID primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Notifications []Notification `json:"notifications" bson:"notifications"`
}

type Notification struct {
	Mail string `json:"mail"`
	Text string `json:"text"`
	CreatedAt string `json:"createdAt"`
	IsOpened bool `json:"isOpened"`
}

type UserNotificationDTO struct {
	ID string `json:"id"`
	Notifications []Notification `json:"notifications"`
}
