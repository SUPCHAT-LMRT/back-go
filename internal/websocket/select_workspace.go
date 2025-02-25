package websocket

import (
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type SelectWorkspace struct {
	DefaultMessage
	WorkspaceId entity.WorkspaceId `json:"workspaceId"`
}

type UnselectWorkspace struct {
	DefaultMessage
}

func (m SelectWorkspace) GetActionName() string {
	return InboundSelectWorkspaceAction
}

func (m UnselectWorkspace) GetActionName() string {
	return InboundUnselectWorkspaceAction
}
