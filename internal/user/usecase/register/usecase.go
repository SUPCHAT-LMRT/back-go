package register

import (
	"context"
	"errors"
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/crypt"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/exists"
	entity2 "github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	delete2 "github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/delete"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/get_data_token_invite"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
	"time"
)

var (
	UserAlreadyExistsErr = errors.New("an account with this email already exists")
)

type RegisterUserDeps struct {
	uberdig.In
	ExistsUserUseCase        *exists.ExistsUserUseCase
	CryptStrategy            crypt.CryptStrategy
	Repository               repository.UserRepository
	Observers                []RegisterUserObserver `group:"register_user_observers"`
	DeleteInviteLinkUseCase  *delete2.DeleteInviteLinkUseCase
	GetInviteLinkDataUseCase *get_data_token_invite.GetInviteLinkDataUseCase
}

type RegisterUserUseCase struct {
	deps RegisterUserDeps
}

func NewRegisterUserUseCase(deps RegisterUserDeps) *RegisterUserUseCase {
	return &RegisterUserUseCase{deps: deps}
}

func (r *RegisterUserUseCase) Execute(ctx context.Context, request RegisterUserRequest) error {

	inviteLinkData, err := r.deps.GetInviteLinkDataUseCase.GetInviteLinkData(ctx, request.Token)
	if err != nil {
		return fmt.Errorf("error getting invite link data: %w", err)
	}

	userExists, err := r.deps.ExistsUserUseCase.Execute(ctx, inviteLinkData.Email)
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

	user := r.EntityUser(request, inviteLinkData)
	user.Id = entity.UserId(bson.NewObjectID().Hex())
	user.CreatedAt = time.Now()

	err = r.deps.Repository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error adding user: %w", err)
	}

	err = r.deps.DeleteInviteLinkUseCase.Execute(ctx, request.Token)
	if err != nil {
		return fmt.Errorf("error deleting invite link: %w", err)
	}

	for _, observer := range r.deps.Observers {
		go observer.NotifyUserRegistered(*user)
	}

	return nil
}

func (r *RegisterUserUseCase) EntityUser(user RegisterUserRequest, link *entity2.InviteLink) *entity.User {
	return &entity.User{
		FirstName: link.FirstName,
		LastName:  link.LastName,
		Email:     link.Email,
		Password:  user.Password,
	}
}

type RegisterUserRequest struct {
	Token    string
	Password string
}
