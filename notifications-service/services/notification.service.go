package services

import "notifications-service/repository"

type NotificationService struct {
	repo *repository.NotificationRepository
}


func NewNotificationService(repo *repository.NotificationRepository) *NotificationService {
	return &NotificationService{
		repo: repo,
	}
}