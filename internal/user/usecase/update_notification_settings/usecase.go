// internal/user/usecase/update_notification_settings/usecase.go
package update_notification_settings

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
	uberdig "go.uber.org/dig"
)

type UpdateNotificationSettingsUseCaseDeps struct {
	uberdig.In
	UserRepository repository.UserRepository
}

type UpdateNotificationSettingsUseCase struct {
	deps UpdateNotificationSettingsUseCaseDeps
}

type UpdateNotificationSettingsRequest struct {
	UserId  entity.UserId
	Enabled bool
}

func NewUpdateNotificationSettingsUseCase(deps UpdateNotificationSettingsUseCaseDeps) *UpdateNotificationSettingsUseCase {
	return &UpdateNotificationSettingsUseCase{deps: deps}
}

func (uc *UpdateNotificationSettingsUseCase) Execute(ctx context.Context, req UpdateNotificationSettingsRequest) error {
	return uc.deps.UserRepository.UpdateNotificationSettings(ctx, req.UserId, req.Enabled)
}
