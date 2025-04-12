package save_status

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/repository"
	uberdig "go.uber.org/dig"
)

type SaveStatusUseCaseDeps struct {
	uberdig.In
	UserStatusRepository repository.UserStatusRepository
	Observers            []UserStatusSavedObserver `group:"save_user_status_observers"`
}

type SaveStatusUseCase struct {
	deps SaveStatusUseCaseDeps
}

func NewSaveStatusUseCase(deps SaveStatusUseCaseDeps) *SaveStatusUseCase {
	return &SaveStatusUseCase{deps: deps}
}

func (u *SaveStatusUseCase) Execute(ctx context.Context, userId user_entity.UserId, status entity.Status) error {
	userStatus := &entity.UserStatus{UserId: userId, Status: status}
	err := u.deps.UserStatusRepository.Save(ctx, userStatus)
	if err != nil {
		return err
	}

	for _, observer := range u.deps.Observers {
		observer.NotifyUserStatusSaved(userStatus)
	}

	return nil
}
