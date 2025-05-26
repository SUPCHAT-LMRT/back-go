package delete

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/repository"
)

type DeleteInviteLinkWorkspaceUseCase struct {
	repository repository.InviteLinkRepository
}

func NewDeleteInviteLinkWorkspaceUseCase(
	linkRepository repository.InviteLinkRepository,
) *DeleteInviteLinkWorkspaceUseCase {
	return &DeleteInviteLinkWorkspaceUseCase{repository: linkRepository}
}

func (d *DeleteInviteLinkWorkspaceUseCase) Execute(ctx context.Context, token string) error {
	return d.repository.DeleteInviteLink(ctx, token)
}
