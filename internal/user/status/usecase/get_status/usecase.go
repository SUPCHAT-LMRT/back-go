package get_status

import (
	"context"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/repository"
)

type GetStatusUseCase struct {
	userStatusRepository repository.UserStatusRepository
}

func NewGetStatusUseCase(userStatusRepository repository.UserStatusRepository) *GetStatusUseCase {
	return &GetStatusUseCase{userStatusRepository: userStatusRepository}
}

func (u *GetStatusUseCase) Execute(
	ctx context.Context,
	userId user_entity.UserId,
) (entity.Status, error) {
	status, err := u.userStatusRepository.Get(ctx, userId)
	if err != nil {
		return entity.StatusUnknown, err
	}

	return status.Status, nil
}
