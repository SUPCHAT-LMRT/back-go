package repository

import (
	"context"
	"errors"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

var ErrUserNotFound = errors.New("user not found")

type UserRepository interface {
	Create(ctx context.Context, user *entity.User) error
	GetById(ctx context.Context, userId entity.UserId) (user *entity.User, err error)
	GetByEmail(
		ctx context.Context,
		userEmail string,
		options ...GetUserOptionFunc,
	) (user *entity.User, err error)
	List(ctx context.Context) (users []*entity.User, err error)
	Update(ctx context.Context, user *entity.User) error
	Delete(ctx context.Context, userId entity.UserId) error
	UpdateNotificationSettings(ctx context.Context, userId entity.UserId, enabled bool) error
}

type getUserOptions struct {
	withPassword bool
}

type GetUserOptionFunc func(options *getUserOptions)

func WithUserPassword() GetUserOptionFunc {
	return func(options *getUserOptions) {
		options.withPassword = true
	}
}
