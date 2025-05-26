package update_banner

import workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"

type SaveBannerWorkspaceObserver interface {
	NotifyUpdateBannerWorkspace(workspace *workspace_entity.Workspace)
}
