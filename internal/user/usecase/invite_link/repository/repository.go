package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	"time"
)

var (
	RedisInviteLinkExpiredTime = 7 * 24 * time.Hour
)

type InviteLinkRepository interface {
	GenerateInviteLink(ctx context.Context, link *entity.InviteLink) error
	GetInviteLinkData(ctx context.Context, token string) (*entity.InviteLink, error)
	DeleteInviteLink(ctx context.Context, token string) error
}
