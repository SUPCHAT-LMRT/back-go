package update_password

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/crypt"
)

type ChangePasswordUseCase struct {
	repository    repository.UserRepository
	cryptStrategy crypt.CryptStrategy
}

func NewChangePasswordUseCase(repository repository.UserRepository, cryptStrategy crypt.CryptStrategy) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{
		repository:    repository,
		cryptStrategy: cryptStrategy,
	}
}

func (u *ChangePasswordUseCase) Execute(ctx context.Context, user *entity.User, password string) error {
	hashedPassword, err := u.cryptStrategy.Hash(password)
	if err != nil {
		panic(err)
	}

	user.Password = hashedPassword
	return u.repository.Update(ctx, user)
}
