package repository

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	"time"
)

var (
	RedisInviteLinkExpiredTime = 7 * 24 * time.Hour
)

var (
	InviteLinkNotFoundErr = errors.New("invite link not found")
)

type InviteLinkRepository interface {
	GenerateInviteLink(ctx context.Context, link *entity.InviteLink) error
	GetInviteLinkData(ctx context.Context, token string) (*entity.InviteLink, error)
	GetInviteLinkDataByEmail(ctx context.Context, email string) (*entity.InviteLink, error)
	DeleteInviteLink(ctx context.Context, token string) error
	GetAllInviteLinks(ctx context.Context) ([]*entity.InviteLink, error)
}
