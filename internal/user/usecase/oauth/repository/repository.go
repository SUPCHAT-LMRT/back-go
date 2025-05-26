package repository

import (
	"context"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/oauth/entity"
)

type OauthConnectionRepository interface {
	CreateOauthConnection(ctx context.Context, connection *entity.OauthConnection) error
	GetOauthConnectionByUserId(ctx context.Context, userId string) (*entity.OauthConnection, error)
	ListOauthConnectionsByUser(
		ctx context.Context,
		userId user_entity.UserId,
	) ([]*entity.OauthConnection, error)
}
