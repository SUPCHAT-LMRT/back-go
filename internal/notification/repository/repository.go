package repository

import (
	"context"
	"errors"

	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

var NotificationNotFoundErr = errors.New("notification not found")

type NotificationRepository interface {
	Create(ctx context.Context, notification *entity.Notification) error
	GetById(ctx context.Context, notificationId entity.NotificationId) (*entity.Notification, error)
	List(ctx context.Context, userId user_entity.UserId) ([]*entity.Notification, error)
	Update(ctx context.Context, notification *entity.Notification) error
	Delete(ctx context.Context, notificationId entity.NotificationId) error
}
