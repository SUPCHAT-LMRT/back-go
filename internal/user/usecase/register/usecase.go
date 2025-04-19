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
	"github.com/supchat-lmrt/back-go/internal/user/usecase/exists_by_oauthemail"
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
	ExistsUserUseCase             *exists_by_email.ExistsUserByEmailUseCase
	ExistsUserByOauthEmailUseCase *exists_by_oauthemail.ExistsUserByOauthEmailUseCase
	CryptStrategy                 crypt.CryptStrategy
	Repository                    repository.UserRepository
	Observers                     []RegisterUserObserver `group:"register_user_observers"`
	DeleteInviteLinkUseCase       *delete2.DeleteInviteLinkUseCase
	GetInviteLinkDataUseCase      *get_data_token_invite.GetInviteLinkDataUseCase
	SearchUserSyncManager         user_search.SearchUserSyncManager
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

	if request.Password != "" {
		hash, err := r.deps.CryptStrategy.Hash(request.Password)
		if err != nil {
			return fmt.Errorf("error hashing password: %w", err)
		}

		request.Password = hash
	} else {
		// If no password is provided, we need to check if the user is using OAuth
		if request.OauthEmail == "" {
			return fmt.Errorf("error: no password provided and no oauth email in invite link")
		}

		// Check if a user with this oauth email already exists
		oauthUserExists, err := r.deps.ExistsUserByOauthEmailUseCase.Execute(ctx, request.OauthEmail)
		if err != nil {
			return fmt.Errorf("error checking if user exists by oauth email: %w", err)
		}
		if oauthUserExists {
			return UserAlreadyExistsErr
		}
	}

	user := r.EntityUser(request, inviteLinkData)
	user.Id = entity.UserId(bson.NewObjectID().Hex())
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt

	err = r.deps.Repository.Create(ctx, user)
	if err != nil {
		return fmt.Errorf("error adding user: %w", err)
	}

	err = r.deps.DeleteInviteLinkUseCase.Execute(ctx, request.Token)
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

	for _, observer := range r.deps.Observers {
		go observer.NotifyUserRegistered(*user)
	}

	return nil
}

func (r *RegisterUserUseCase) EntityUser(user RegisterUserRequest, link *entity2.InviteLink) *entity.User {
	return &entity.User{
		FirstName:  link.FirstName,
		LastName:   link.LastName,
		Email:      link.Email,
		OauthEmail: user.OauthEmail,
		Password:   user.Password,
	}
}

type RegisterUserRequest struct {
	Token string
	// OauthEmail is the email used for OAuth registration (optional if password is provided)
	OauthEmail string
	// Password is the password used for registration (optional if oauth)
	Password string
}
