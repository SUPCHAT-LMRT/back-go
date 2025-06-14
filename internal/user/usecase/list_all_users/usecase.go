package list_all_users

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type ListUserUseCase struct {
	repo repository.UserRepository
}

func NewListUserUseCase(repo repository.UserRepository) *ListUserUseCase {
	return &ListUserUseCase{repo: repo}
}

func (u *ListUserUseCase) Execute(ctx context.Context) ([]*entity.User, error) {
	return u.repo.List(ctx)
}
