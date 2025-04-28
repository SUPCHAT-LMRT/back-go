package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type OutboundChannelsReordered struct {
	messages.DefaultMessage
	ChannelReorders []ChannelReorderMessage `json:"channelReorders"`
}

type OutboundChannelsDeleted struct {
	messages.DefaultMessage
	ChannelId   entity.ChannelId             `json:"channelId"`
	WorkspaceId workspace_entity.WorkspaceId `json:"workspaceId"`
}

type ChannelReorderMessage struct {
	ChannelId entity.ChannelId `json:"channelId"`
	NewOrder  int              `json:"newOrder"`
}

func (m OutboundChannelsDeleted) GetActionName() messages.Action {
	return messages.OutboundChannelsDeletedAction
}

func (m OutboundChannelsDeleted) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}

func (m OutboundChannelsReordered) GetActionName() messages.Action {
	return messages.OutboundChannelsReorderedAction
}

func (m OutboundChannelsReordered) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
