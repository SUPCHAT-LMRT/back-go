package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
)

type InviteLinkRepository interface {
	GenerateInviteLink(ctx context.Context, link *entity.InviteLink) error
	GetInviteLinkData(ctx context.Context, token string) (*entity.InviteLink, error)
}
