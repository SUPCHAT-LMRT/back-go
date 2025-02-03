package update_icon

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"io"
)

type UpdateWorkspaceIconStrategy interface {
	Handle(ctx context.Context, workspaceId entity.WorkspaceId, imageReader io.Reader, contentType string) error
}
