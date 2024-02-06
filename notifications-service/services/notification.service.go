package services

import (
	"context"
	"fmt"
	"notifications-service/client"
	"notifications-service/config"
	"notifications-service/domains"
	"notifications-service/errors"
	"notifications-service/repository"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.opentelemetry.io/otel/trace"
)

type NotificationService struct {
	repo       *repository.NotificationRepository
	mailClient *client.MailClient
	userClient *client.UserClient
	logger     *config.Logger
	tracer     trace.Tracer
}

func NewNotificationService(repo *repository.NotificationRepository, mailClient *client.MailClient, userClient *client.UserClient, logger *config.Logger, tracer trace.Tracer) *NotificationService {
	return &NotificationService{
		repo:       repo,
		mailClient: mailClient,
		userClient: userClient,
		logger:     logger,
		tracer:     tracer,
	}
}

func (ns NotificationService) CreateNewUserNotification(ctx context.Context, id string) (*domains.UserNotificationDTO, *errors.ErrorStruct) {
	ctx, span := ns.tracer.Start(ctx, "NotificationService.CreateNewUserNotification")
	defer span.End()
	ns.logger.LogInfo("notification-service", fmt.Sprintf("Trying to create a user structure for user with id %v", id))
	userNotification, err := ns.repo.CreateNewUserNotification(ctx, id)
	if err != nil {
		ns.logger.LogError("notification-service", fmt.Sprintf("Error while creating a user structure for user with id %v", id))
		return nil, err
	}
	ns.logger.LogInfo("notification-service", fmt.Sprintf("Successfully created a user structure for user with id %v", id))
	return &domains.UserNotificationDTO{
		ID:            userNotification.ID.Hex(),
		Notifications: userNotification.Notifications,
	}, nil
}

func (ns NotificationService) PushNewNotificationToUser(ctx context.Context, id string, notification domains.Notification) (*domains.UserNotificationDTO, *errors.ErrorStruct) {
	ctx, span := ns.tracer.Start(ctx, "NotificationService.PushNewNotificationToUser")
	defer span.End()
	ns.logger.LogInfo("notification-service", fmt.Sprintf("Trying to psuh notification for user with id %v", id))
	userNotification, err := ns.repo.FindOneUserNotificationByID(ctx, id)
	if err != nil {
		ns.logger.LogError("notification-service", err.GetErrorMessage())
		return nil, err
	}
	notification.CreatedAt = time.Now().String()
	notification.IsOpened = false
	userNotification.Notifications = append([]domains.Notification{notification}, userNotification.Notifications...)
	repoResp, err := ns.repo.UpdateNotificationByID(ctx, userNotification)
	if err != nil {
		ns.logger.LogError("notification-service", err.GetErrorMessage())
		return nil, err
	}
	user, errFromUser := ns.userClient.GetAllInformationsByUserID(ctx, id)
	if errFromUser != nil {
		ns.logger.LogError("notification-service", err.GetErrorMessage())
		return nil, errFromUser
	}
	go func() {
		ns.mailClient.SendMailNotification(notification, user.Email)
	}()
	ns.logger.LogInfo("notification-service", fmt.Sprintf("Successfully pushed new notiricationj %v", id))
	return &domains.UserNotificationDTO{
		ID:            repoResp.ID.Hex(),
		Notifications: repoResp.Notifications,
	}, nil
}

func (ns NotificationService) ReadAllNotifications(ctx context.Context, notifications domains.UserNotificationDTO) (*domains.UserNotificationDTO, *errors.ErrorStruct) {
	ctx, span := ns.tracer.Start(ctx, "NotificationService.ReadAllNotifications")
	defer span.End()
	ns.logger.LogInfo("notification-service", fmt.Sprintf("Trying to read all notifications for user with id %v", notifications.ID))
	notifications.Notifications = *ns.makeAllNotificationsOpened(&notifications.Notifications)
	castedKey, _ := primitive.ObjectIDFromHex(notifications.ID)
	castedUserNotification := domains.UserNotification{
		ID:            castedKey,
		Notifications: notifications.Notifications,
	}
	_, err := ns.repo.UpdateNotificationByID(ctx, &castedUserNotification)
	if err != nil {
		ns.logger.LogError("notification-service", err.GetErrorMessage())
		return nil, err
	}
	ns.logger.LogInfo("notification-service", fmt.Sprintf("Readed all notifications for id %v", notifications.ID))
	return &notifications, nil
}

func (ns NotificationService) FindAllNotificationsByID(ctx context.Context, id string) (*domains.UserNotificationDTO, *errors.ErrorStruct) {
	ctx, span := ns.tracer.Start(ctx, "NotificationService.FindAllNotificationsByID")
	defer span.End()
	notifications, err := ns.repo.FindOneUserNotificationByID(ctx, id)
	if err != nil {
		ns.logger.LogError("notification-service", err.GetErrorMessage())
		return nil, err
	}
	ns.logger.LogInfo("notification-service", fmt.Sprintf("Readed all notifications for id %v", id))
	return &domains.UserNotificationDTO{
		ID:            notifications.ID.Hex(),
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
