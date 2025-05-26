package delete_user

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type DeleteUserUseCase struct {
	repository repository.UserRepository
}

func NewDeleteUserUseCase(repo repository.UserRepository) *DeleteUserUseCase {
	return &DeleteUserUseCase{repository: repo}
}

func (uc *DeleteUserUseCase) Execute(ctx context.Context, userId entity.UserId) error {
	err := uc.repository.Delete(ctx, userId)
	if err != nil {
		return err
	}
	return nil
}
