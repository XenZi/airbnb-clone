package services

import (
	"context"
	"log"
	"notifications-service/client"
	"notifications-service/domains"
	"notifications-service/errors"
	"notifications-service/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type NotificationService struct {
	repo *repository.NotificationRepository
	mailClient *client.MailClient
	userClient *client.UserClient
}


func NewNotificationService(repo *repository.NotificationRepository, mailClient *client.MailClient, userClient *client.UserClient) *NotificationService {
	return &NotificationService{
		repo: repo,
		mailClient: mailClient,
		userClient: userClient,
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

func (ns NotificationService) PushNewNotificationToUser(ctx context.Context, id string, notification domains.Notification) (*domains.UserNotificationDTO, *errors.ErrorStruct){
	userNotification, err := ns.repo.FindOneUserNotificationByID(id)
	if err != nil {
		return nil, err
	}
	notification.CreatedAt = time.Now().String()
	notification.IsOpened = false
	userNotification.Notifications = append([]domains.Notification{notification}, userNotification.Notifications...)
	repoResp, err := ns.repo.UpdateNotificationByID(userNotification)
	if err != nil {
		return nil, err
	}
	user, errFromUser := ns.userClient.GetAllInformationsByUserID(ctx, id)
	if errFromUser != nil {
		return nil, errFromUser
	}
    go func(ctx context.Context) {
        // Check if the context is done before proceeding
        select {
        case <-ctx.Done():
            // Log the context cancellation or timeout
            log.Printf("Context cancelled or timed out: %v", ctx.Err())
            return
        default:
            // Context is not done, proceed with sending mail
            if err := ns.mailClient.SendMailNotification(ctx, notification, user.Email); err != nil {
                // Log the error
                log.Printf("Error sending mail notification: %v", err)
            }
        }
    }(ctx)

	return &domains.UserNotificationDTO{
		ID: repoResp.ID.Hex(),
		Notifications: repoResp.Notifications,
	}, nil
}

func (ns NotificationService) ReadAllNotifications(notifications domains.UserNotificationDTO) (*domains.UserNotificationDTO, *errors.ErrorStruct) {
	notifications.Notifications = *ns.makeAllNotificationsOpened(&notifications.Notifications)
	castedKey, _ := primitive.ObjectIDFromHex(notifications.ID)
	castedUserNotification := domains.UserNotification{
		ID: castedKey,
		Notifications: notifications.Notifications,
	}
	_, err := ns.repo.UpdateNotificationByID(&castedUserNotification)
	if err != nil {
		return nil, err
	}
	return &notifications, nil
}

func (ns NotificationService) FindAllNotificationsByID(id string) (*domains.UserNotificationDTO, *errors.ErrorStruct) {
	notifications, err := ns.repo.FindOneUserNotificationByID(id)
	if err != nil {
		return nil, err
	}
	return &domains.UserNotificationDTO{
		ID: notifications.ID.Hex(),
		Notifications: notifications.Notifications,
	}, nil
}

func (ns NotificationService) makeAllNotificationsOpened(notifications *[]domains.Notification) *[]domains.Notification {
    for i := range *notifications {
		if (*notifications)[i].IsOpened {
			break
		}
		(*notifications)[i].IsOpened = true
	}
    return notifications
}