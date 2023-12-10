package services

import (
	"notifications-service/domains"
	"notifications-service/errors"
	"notifications-service/repository"
	"time"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}


func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}

func (ns NotificationService) CreateNewUserNotification(id string)  (*domains.UserNotificationDTO, *errors.ErrorStruct) {
	userNotification, err := ns.repo.CreateNewUserNotification(id)
	if err != nil {
		return nil, err
	}
	return &domains.UserNotificationDTO{
		ID: userNotification.ID.Hex(),
		Notifications: userNotification.Notifications,
	}, nil
}

func (ns NotificationService) PushNewNotificationToUser(id string, notification domains.Notification) (*domains.UserNotificationDTO, *errors.ErrorStruct){
	userNotification, err := ns.repo.FindOneUserNotificationByID(id)
	if err != nil {
		return nil, err
	}
	notification.CreatedAt = time.Now().String()
	notification.IsOpened = false
	userNotification.Notifications = append(userNotification.Notifications, notification)
	repoResp, err := ns.repo.UpdateNotificationByID(userNotification)
	if err != nil {
		return nil, err
	}
	return &domains.UserNotificationDTO{
		ID: repoResp.ID.Hex(),
		Notifications: repoResp.Notifications,
	}, nil
}