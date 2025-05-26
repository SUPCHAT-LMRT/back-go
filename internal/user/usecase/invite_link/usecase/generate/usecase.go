package generate

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
	user_repository "github.com/supchat-lmrt/back-go/internal/user/repository"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_email"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/repository"
	uberdig "go.uber.org/dig"
)

type InviteLinkUseCaseDeps struct {
	uberdig.In
	Repository            repository.InviteLinkRepository
	GetUserByEmailUseCase *get_by_email.GetUserByEmailUseCase
	Observers             []GenerateInviteLinkObserver `group:"generate_invite_link_observers"`
}

type InviteLinkUseCase struct {
	deps InviteLinkUseCaseDeps
}

func NewInviteLinkUseCase(deps InviteLinkUseCaseDeps) *InviteLinkUseCase {
	return &InviteLinkUseCase{deps: deps}
}

func (u *InviteLinkUseCase) CreateInviteLink(
	ctx context.Context,
	firstName, lastName, email string,
) (string, error) {
	_, err := u.deps.GetUserByEmailUseCase.Execute(ctx, email)
	if err == nil {
		return "", errors.New("user already exists")
	} else if !errors.Is(err, user_repository.ErrUserNotFound) {
		return "", err
	}

	_, err = u.deps.Repository.GetInviteLinkDataByEmail(ctx, email)
	if err == nil {
		return "", errors.New("user already invited")
	} else if !errors.Is(err, repository.ErrInviteLinkNotFound) {
		return "", err
	}

	token := uuid.New().String()
	inviteLink := &entity.InviteLink{
		Token:     token,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	err = u.deps.Repository.GenerateInviteLink(ctx, inviteLink)
	if err != nil {
		return "", err
	}

	inviteLinkFormat := os.Getenv("FRONT_ACCOUNT_REGISTER_URL")
	if inviteLinkFormat = strings.Replace(inviteLinkFormat, "{token}", inviteLink.Token, 1); inviteLinkFormat == "" {
		return "", errors.New("invite link format is empty")
	}

	for _, observer := range u.deps.Observers {
		observer.NotifyInviteLinkGenerated(inviteLink, inviteLinkFormat)
	}

	return inviteLinkFormat, nil
}
