package repository

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
)

type UserStatusRepository interface {
	Get(ctx context.Context, userId user_entity.UserId) (*entity.UserStatus, error)
	Save(ctx context.Context, userStatus *entity.UserStatus) error
}
