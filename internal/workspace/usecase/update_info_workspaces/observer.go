package update_info_workspaces

import (
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type SaveInfoWorkspacesObserver interface {
	NotifyUpdateInfoWorkspaces(workspace *workspace_entity.Workspace)
}
