package update_banner

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"io"
)

type UpdateWorkspaceBannerStrategy interface {
	Handle(ctx context.Context, workspaceId entity.WorkspaceId, imageReader io.Reader, contentType string) error
}
