package update_user

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type UpdateUserUseCase struct {
	repository repository.UserRepository
}

func NewUpdateUserUseCase(repository repository.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repository: repository}
}

func (u *UpdateUserUseCase) Execute(ctx context.Context, user *user_entity.User) error {
	return u.repository.Update(ctx, user)
}
