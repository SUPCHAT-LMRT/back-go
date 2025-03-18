package generate

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/repository"
	"os"
	"strings"
)

type InviteLinkUseCase struct {
	repository repository.InviteLinkRepository
}

func NewInviteLinkUseCase(linkRepository repository.InviteLinkRepository) *InviteLinkUseCase {
	return &InviteLinkUseCase{repository: linkRepository}
}

func (u *InviteLinkUseCase) CreateInviteLink(ctx context.Context, firstName, lastName, email string) (string, error) {
	token := uuid.New().String()
	inviteLink := &entity.InviteLink{
		Token:     token,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}

	err := u.repository.GenerateInviteLink(ctx, inviteLink)
	if err != nil {
		return "", err
	}

	inviteLinkFormat := os.Getenv("FRONT_ACCOUNT_REGISTER_URL")
	if inviteLinkFormat = strings.Replace(inviteLinkFormat, "{token}", inviteLink.Token, 1); inviteLinkFormat == "" {
		return "", errors.New("invite link format is empty")
	}

	return inviteLinkFormat, nil
}
