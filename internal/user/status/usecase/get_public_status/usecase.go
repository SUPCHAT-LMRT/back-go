package get_public_status

import (
	"context"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_or_create_status"
)

type GetPublicStatusUseCase struct {
	useCase *get_or_create_status.GetOrCreateStatusUseCase
}

func NewGetPublicStatusUseCase(
	useCase *get_or_create_status.GetOrCreateStatusUseCase,
) *GetPublicStatusUseCase {
	return &GetPublicStatusUseCase{useCase: useCase}
}

func (u *GetPublicStatusUseCase) Execute(
	ctx context.Context,
	userId user_entity.UserId,
	defaultStatus entity.Status,
) (entity.Status, error) {
	status, err := u.useCase.Execute(ctx, userId, defaultStatus)
	if err != nil {
		return entity.StatusUnknown, err
	}

	if status == entity.StatusInvisible {
		status = entity.StatusOffline
	}

	return status, nil
}
