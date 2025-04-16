package event

import workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"

const (
	WorkspaceUpdatedEventType EventType = "workspace_updated"
)

type WorkspaceUpdatedEvent struct {
	Workspace *workspace_entity.Workspace
}

func (e WorkspaceUpdatedEvent) Type() EventType {
	return WorkspaceUpdatedEventType
}
