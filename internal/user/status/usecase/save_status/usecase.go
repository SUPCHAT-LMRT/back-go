package save_status

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/repository"
)

type SaveStatusUseCase struct {
	userStatusRepository repository.UserStatusRepository
}

func NewSaveStatusUseCase(userStatusRepository repository.UserStatusRepository) *SaveStatusUseCase {
	return &SaveStatusUseCase{userStatusRepository: userStatusRepository}
}

func (u *SaveStatusUseCase) Execute(ctx context.Context, userId user_entity.UserId, status entity.Status) error {
	return u.userStatusRepository.Save(ctx, &entity.UserStatus{UserId: userId, Status: status})
}
