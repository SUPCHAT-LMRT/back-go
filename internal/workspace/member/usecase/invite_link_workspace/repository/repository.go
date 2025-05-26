package repository

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/entity"
)

var RedisInviteLinkExpiredTime = 24 * time.Hour

type InviteLinkRepository interface {
	GenerateInviteLink(ctx context.Context, link *entity.InviteLink) error
	GetInviteLinkData(ctx context.Context, token string) (*entity.InviteLink, error)
	DeleteInviteLink(ctx context.Context, token string) error
}
