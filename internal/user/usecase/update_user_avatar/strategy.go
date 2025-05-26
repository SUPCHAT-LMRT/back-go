package update_user_avatar

import (
	"context"
	"io"

	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type UpdateUserAvatarStrategy interface {
	Handle(
		ctx context.Context,
		userId entity.UserId,
		imageReader io.Reader,
		contentType string,
	) error
}
