package inbound

import (
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type InboundSelectWorkspace struct {
	messages.DefaultMessage
	WorkspaceId workspace_entity.WorkspaceId `json:"workspaceId"`
}

type InboundUnselectWorkspace struct {
	messages.DefaultMessage
}

func (m InboundSelectWorkspace) GetActionName() messages.Action {
	return messages.InboundSelectWorkspaceAction
}

func (m InboundUnselectWorkspace) GetActionName() messages.Action {
	return messages.InboundUnselectWorkspaceAction
}
