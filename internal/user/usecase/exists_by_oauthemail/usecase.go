package exists_by_oauthemail

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
)

type ExistsUserByOauthEmailUseCase struct {
	userRepository repository.UserRepository
}

func NewExistsUserByOauthEmailUseCase(userRepository repository.UserRepository) *ExistsUserByOauthEmailUseCase {
	return &ExistsUserByOauthEmailUseCase{userRepository: userRepository}
}

func (u *ExistsUserByOauthEmailUseCase) Execute(ctx context.Context, oauthEmail string) (bool, error) {
	_, err := u.userRepository.GetByOauthEmail(ctx, oauthEmail)
	if err != nil {
		if errors.Is(err, repository.UserNotFoundErr) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
