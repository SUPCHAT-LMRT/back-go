package update_icon

import workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"

type SaveIconWorkspaceObserver interface {
	NotifyUpdateBannerWorkspace(workspace *workspace_entity.Workspace)
}
