package delete

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/repository"
)

type DeleteInviteLinkUseCase struct {
	repository repository.InviteLinkRepository
}

func NewDeleteInviteLinkUseCase(linkRepository repository.InviteLinkRepository) *DeleteInviteLinkUseCase {
	return &DeleteInviteLinkUseCase{repository: linkRepository}
}

func (d *DeleteInviteLinkUseCase) Execute(ctx context.Context, token string) error {
	return d.repository.DeleteInviteLink(ctx, token)
}
