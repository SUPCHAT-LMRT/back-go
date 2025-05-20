package update_banner

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	uberdig "go.uber.org/dig"
)

type UpdateWorkspaceIconDeps struct {
	uberdig.In
	EventBus *event.EventBus
	Logger   logger.Logger
}

type UpdateWorkspaceIconObserver struct {
	deps UpdateWorkspaceIconDeps
}

func NewUpdateWorkspaceBannerObserver(deps UpdateWorkspaceIconDeps) SaveBannerWorkspaceObserver {
	return &UpdateWorkspaceIconObserver{deps: deps}
}

func (o UpdateWorkspaceIconObserver) NotifyUpdateBannerWorkspace(workspaces *entity.Workspace) {
	o.deps.EventBus.Publish(&event.WorkspaceUpdatedEvent{
		Workspace: workspaces,
	})
}
