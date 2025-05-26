package mark_as_read

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	"github.com/supchat-lmrt/back-go/internal/notification/repository"
)

type MarkAsReadUseCase struct {
	notificationRepository repository.NotificationRepository
}

func NewMarkAsReadUseCase(
	notificationRepository repository.NotificationRepository,
) *MarkAsReadUseCase {
	return &MarkAsReadUseCase{notificationRepository: notificationRepository}
}

func (u *MarkAsReadUseCase) Execute(
	ctx context.Context,
	notificationId entity.NotificationId,
) error {
	notification, err := u.notificationRepository.GetById(ctx, notificationId)
	if err != nil {
		return err
	}

	notification.IsRead = true
	return u.notificationRepository.Update(ctx, notification)
}
