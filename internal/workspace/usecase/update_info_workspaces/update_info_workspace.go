package update_info_workspaces

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	uberdig "go.uber.org/dig"
)

type UpdateInfoWorkspacesDeps struct {
	uberdig.In
	EventBus *event.EventBus
	Logger   logger.Logger
}

type UpdateInfoWorkspacesObserver struct {
	deps UpdateInfoWorkspacesDeps
}

func NewUpdateInfoWorkspacesObserver(deps UpdateInfoWorkspacesDeps) SaveInfoWorkspacesObserver {
	return &UpdateInfoWorkspacesObserver{deps: deps}
}

func (o UpdateInfoWorkspacesObserver) NotifyUpdateInfoWorkspaces(workspaces *entity.Workspace) {
	o.deps.EventBus.Publish(&event.WorkspaceUpdatedEvent{
		Workspace: workspaces,
	})
}
