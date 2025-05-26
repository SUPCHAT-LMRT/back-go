package exists_by_email

import (
	"context"
	"errors"

	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type ExistsUserByEmailUseCase struct {
	userRepository repository.UserRepository
}

func NewExistsUserByEmailUseCase(
	userRepository repository.UserRepository,
) *ExistsUserByEmailUseCase {
	return &ExistsUserByEmailUseCase{userRepository: userRepository}
}

func (u *ExistsUserByEmailUseCase) Execute(ctx context.Context, userEmail string) (bool, error) {
	_, err := u.userRepository.GetByEmail(ctx, userEmail)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
