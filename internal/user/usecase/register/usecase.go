package register

import (
	"context"
	"errors"
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/crypt"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/exists"
	uberdig "go.uber.org/dig"
)

var (
	UserAlreadyExistsErr = errors.New("an account with this email already exists")
)

type RegisterUserDeps struct {
	uberdig.In
	ExistsUserUseCase *exists.ExistsUserUseCase
	CryptStrategy     crypt.CryptStrategy
	Repository        repository.UserRepository
	Observers         []RegisterUserObserver `group:"register_user_observers"`
}

type RegisterUserUseCase struct {
	deps RegisterUserDeps
}

func NewRegisterUserUseCase(deps RegisterUserDeps) *RegisterUserUseCase {
	return &RegisterUserUseCase{deps: deps}
}

func (r *RegisterUserUseCase) Execute(ctx context.Context, request RegisterUserRequest) error {

	userExists, err := r.deps.ExistsUserUseCase.Execute(ctx, request.Email)
	if err != nil {
		return fmt.Errorf("error checking if user exists: %w", err)
	}
	if userExists {
		return UserAlreadyExistsErr
	}

	hash, err := r.deps.CryptStrategy.Hash(request.Password)
	if err != nil {
		return fmt.Errorf("error hashing password: %w", err)
	}

	request.Password = hash

	user := r.EntityUser(request)
	err = r.deps.Repository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error adding user: %w", err)
	}

	for _, observer := range r.deps.Observers {
		go observer.NotifyUserRegistered(*user)
	}

	return nil
}

func (r *RegisterUserUseCase) EntityUser(user RegisterUserRequest) *entity.User {
	return &entity.User{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		Password:  user.Password,
	}
}

type RegisterUserRequest struct {
	FirstName string
	LastName  string
	Email     string
	Password  string
}
