package update_type_workspace

import (
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type SaveTypeWorkspaceObserver interface {
	NotifyUpdateTypeWorkspace(workspace *workspace_entity.Workspace)
}
