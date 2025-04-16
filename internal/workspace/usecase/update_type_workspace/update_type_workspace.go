package update_type_workspace

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	uberdig "go.uber.org/dig"
)

type UpdateTypeWorkspaceDeps struct {
	uberdig.In
	EventBus *event.EventBus
	Logger   logger.Logger
}

type UpdateTypeWorkspaceObserver struct {
	deps UpdateTypeWorkspaceDeps
}

func NewNotifyUpdateTypeWorkspaceObserver(deps UpdateTypeWorkspaceDeps) SaveTypeWorkspaceObserver {
	return &UpdateTypeWorkspaceObserver{deps: deps}
}

func (o UpdateTypeWorkspaceObserver) NotifyUpdateTypeWorkspace(workspaces *entity.Workspace) {
	o.deps.EventBus.Publish(&event.WorkspaceUpdatedEvent{
		Workspace: workspaces,
	})
}
