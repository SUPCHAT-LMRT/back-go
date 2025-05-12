package list_all_users

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type ListUserUseCase interface {
	Execute(ctx context.Context) ([]*entity.User, error)
}

type useCase struct {
	repo repository.UserRepository
}

func NewListUserUseCase(repo repository.UserRepository) ListUserUseCase {
	return &useCase{repo: repo}
}

func (u *useCase) Execute(ctx context.Context) ([]*entity.User, error) {
	return u.repo.List(ctx)
}
