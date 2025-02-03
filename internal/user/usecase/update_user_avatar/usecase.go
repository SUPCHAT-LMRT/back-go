package update_user_avatar

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
	"io"
)

type UpdateUserAvatarUseCaseDeps struct {
	uberdig.In
	Strategy UpdateUserAvatarStrategy
}

type UpdateUserAvatarUseCase struct {
	deps UpdateUserAvatarUseCaseDeps
}

func NewUpdateUserAvatarUseCase(deps UpdateUserAvatarUseCaseDeps) *UpdateUserAvatarUseCase {
	return &UpdateUserAvatarUseCase{deps: deps}
}

func (u *UpdateUserAvatarUseCase) Execute(ctx context.Context, userId entity.UserId, updateAvatar UpdateAvatar) error {
	return u.deps.Strategy.Handle(ctx, userId, updateAvatar.ImageReader, updateAvatar.ContentType)
}

type UpdateAvatar struct {
	ImageReader io.Reader
	ContentType string
}
