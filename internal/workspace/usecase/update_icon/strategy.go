package update_icon

import (
	"context"
	"io"

	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type UpdateWorkspaceIconStrategy interface {
	Handle(
		ctx context.Context,
		workspaceId entity.WorkspaceId,
		imageReader io.Reader,
		contentType string,
	) error
}
