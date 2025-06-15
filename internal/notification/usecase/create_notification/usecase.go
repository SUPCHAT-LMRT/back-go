package create_notification

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	"github.com/supchat-lmrt/back-go/internal/notification/repository"
	uberdig "go.uber.org/dig"
)

type CreateNotificationUseCaseDeps struct {
	uberdig.In
	NotificationRepository repository.NotificationRepository
}

type CreateNotificationUseCase struct {
	deps CreateNotificationUseCaseDeps
}

func NewCreateNotificationUseCase(deps CreateNotificationUseCaseDeps) *CreateNotificationUseCase {
	return &CreateNotificationUseCase{deps: deps}
}

func (uc *CreateNotificationUseCase) Execute(ctx context.Context, notification *entity.Notification) error {
	return uc.deps.NotificationRepository.Create(ctx, notification)
}
