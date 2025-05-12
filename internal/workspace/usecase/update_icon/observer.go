package update_icon

import workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"

type SaveIconWorkspaceObserver interface {
	NotifyUpdateIconWorkspace(workspace *workspace_entity.Workspace)
}
