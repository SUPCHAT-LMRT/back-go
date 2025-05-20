package list_notifications

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	"github.com/supchat-lmrt/back-go/internal/notification/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ListNotificationsUseCase struct {
	notificationRepository repository.NotificationRepository
}

func NewListNotificationsUseCase(notificationRepository repository.NotificationRepository) *ListNotificationsUseCase {
	return &ListNotificationsUseCase{notificationRepository: notificationRepository}
}

func (u *ListNotificationsUseCase) Execute(ctx context.Context, userId user_entity.UserId) ([]*entity.Notification, error) {
	return u.notificationRepository.List(ctx, userId)
}
