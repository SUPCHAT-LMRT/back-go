package get_by_oauthemail

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type GetUserByOauthEmailUseCase struct {
	userRepository repository.UserRepository
}

func NewGetUserByOauthEmailUseCase(userRepository repository.UserRepository) *GetUserByOauthEmailUseCase {
	return &GetUserByOauthEmailUseCase{userRepository: userRepository}
}

func (u *GetUserByOauthEmailUseCase) Execute(ctx context.Context, oauthEmail string) (*entity.User, error) {
	return u.userRepository.GetByOauthEmail(ctx, oauthEmail)
}
