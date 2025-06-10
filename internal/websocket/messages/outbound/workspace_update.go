package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundWorkspaceUpdated struct {
	messages.DefaultMessage
	WorkspaceId string `json:"workspaceId"`
	Name        string `json:"name"`
	Topic       string `json:"topic"`
	Type        string `json:"type"`
}

func (m *OutboundWorkspaceUpdated) GetActionName() messages.Action {
	return messages.OutboundWorkspaceUpdatedAction
}

func (m *OutboundWorkspaceUpdated) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
