package register

import (
	"context"
	"errors"
	"fmt"
	user_search "github.com/supchat-lmrt/back-go/internal/search/user"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/crypt"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/exists_by_email"
	entity2 "github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	delete2 "github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/delete"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/usecase/get_data_token_invite"
	oauth_entity "github.com/supchat-lmrt/back-go/internal/user/usecase/oauth/entity"
	user_oauth_repository "github.com/supchat-lmrt/back-go/internal/user/usecase/oauth/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
	"time"
)

var (
	UserAlreadyExistsErr = errors.New("an account with this email already exists")
)

type RegisterUserDeps struct {
	uberdig.In
	ExistsUserUseCase         *exists_by_email.ExistsUserByEmailUseCase
	CryptStrategy             crypt.CryptStrategy
	Repository                repository.UserRepository
	Observers                 []RegisterUserObserver `group:"register_user_observers"`
	DeleteInviteLinkUseCase   *delete2.DeleteInviteLinkUseCase
	GetInviteLinkDataUseCase  *get_data_token_invite.GetInviteLinkDataUseCase
	SearchUserSyncManager     user_search.SearchUserSyncManager
	OauthConnectionRepository user_oauth_repository.OauthConnectionRepository
}

type RegisterUserUseCase struct {
	deps RegisterUserDeps
}

func NewRegisterUserUseCase(deps RegisterUserDeps) *RegisterUserUseCase {
	return &RegisterUserUseCase{deps: deps}
}

func (r *RegisterUserUseCase) Execute(ctx context.Context, token string, opts ...RegisterOption) error {
	options := RegisterOptions{}
	for _, opt := range opts {
		opt(&options)
	}

	inviteLinkData, err := r.deps.GetInviteLinkDataUseCase.GetInviteLinkData(ctx, token)
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

	if options.Mode == RegisterModePassword {
		hash, err := r.deps.CryptStrategy.Hash(options.Password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}

		options.Password = hash
	}

	user := r.EntityUser(options.Password, inviteLinkData)
	user.Id = entity.UserId(bson.NewObjectID().Hex())
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	err = r.deps.Repository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error adding user: %w", err)
	}

	err = r.deps.DeleteInviteLinkUseCase.Execute(ctx, token)
	if err != nil {
		return fmt.Errorf("error deleting invite link: %w", err)
	}

	err = r.deps.SearchUserSyncManager.AddUser(ctx, &user_search.SearchUser{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
	if err != nil {
		return fmt.Errorf("error syncing user: %w", err)
	}

	if options.Mode == RegisterModeOauth {
		// Handle Oauth binding between user and provider
		err = r.deps.OauthConnectionRepository.CreateOauthConnection(ctx, &oauth_entity.OauthConnection{
			UserId:      user.Id,
			Provider:    options.Oauth.Provider,
			OauthEmail:  options.Oauth.Email,
			OauthUserId: options.Oauth.UserId,
		})
		if err != nil {
			return fmt.Errorf("error creating oauth connection: %w", err)
		}
	}

	for _, observer := range r.deps.Observers {
		go observer.NotifyUserRegistered(*user)
	}

	return nil
}

func (r *RegisterUserUseCase) EntityUser(password string, link *entity2.InviteLink) *entity.User {
	return &entity.User{
		FirstName: link.FirstName,
		LastName:  link.LastName,
		Email:     link.Email,
		Password:  password,
	}
}

type RegisterMode uint

const (
	RegisterModePassword RegisterMode = iota
	RegisterModeOauth
)

type RegisterOptions struct {
	Mode     RegisterMode
	Password string
	Oauth    OauthRegisterOptions
}

type OauthRegisterOptions struct {
	Provider string
	UserId   string
	Email    string
}

type RegisterOption func(*RegisterOptions)

func WithPassword(password string) RegisterOption {
	return func(options *RegisterOptions) {
		options.Password = password
		options.Mode = RegisterModePassword
	}
}

func WithOauth(provider string, userId string, email string) RegisterOption {
	return func(options *RegisterOptions) {
		options.Oauth.Provider = provider
		options.Oauth.UserId = userId
		options.Oauth.Email = email
		options.Mode = RegisterModeOauth
	}
}
