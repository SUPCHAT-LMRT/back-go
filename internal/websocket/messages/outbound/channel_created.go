package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type OutboundChannelCreated struct {
	messages.DefaultMessage
	Channel OutboundChannelCreatedChannel `json:"channel"`
}

type OutboundChannelCreatedChannel struct {
	Id          channel_entity.ChannelId     `json:"id"`
	Name        string                       `json:"name"`
	Topic       string                       `json:"topic"`
	WorkspaceId workspace_entity.WorkspaceId `json:"workspaceId"`
}

func (m OutboundChannelCreated) GetActionName() messages.Action {
	return messages.OutboundChannelCreatedAction
}

func (m OutboundChannelCreated) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
