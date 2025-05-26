package get_by_id

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type GetUserByIdUseCase struct {
	userRepository repository.UserRepository
}

func NewGetUserByIdUseCase(userRepository repository.UserRepository) *GetUserByIdUseCase {
	return &GetUserByIdUseCase{userRepository: userRepository}
}

// TODO impl cache handling
func (u *GetUserByIdUseCase) Execute(
	ctx context.Context,
	userId entity.UserId,
) (*entity.User, error) {
	return u.userRepository.GetById(ctx, userId)
}
