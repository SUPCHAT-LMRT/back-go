package generate

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/google/uuid"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/repository"
)

type InviteLinkUseCase struct {
	repository repository.InviteLinkRepository
}

func NewInviteLinkUseCase(linkRepository repository.InviteLinkRepository) *InviteLinkUseCase {
	return &InviteLinkUseCase{repository: linkRepository}
}

func (u *InviteLinkUseCase) Execute(
	ctx context.Context,
	workspaceId workspace_entity.WorkspaceId,
) (string, error) {
	token := uuid.New().String()

	inviteLink := &entity.InviteLink{
		Token:       token,
		WorkspaceId: workspaceId,
	}
	err := u.repository.GenerateInviteLink(ctx, inviteLink)
	if err != nil {
		return "", err
	}

	inviteLinkFormat := os.Getenv("FRONT_WORKSPACE_INVITE_URL")
	if inviteLinkFormat = strings.Replace(inviteLinkFormat, "{token}", inviteLink.Token, 1); inviteLinkFormat == "" {
		return "", errors.New("invite link format is empty")
	}

	return inviteLinkFormat, nil
}
