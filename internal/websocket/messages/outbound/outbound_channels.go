package outbound

import (
	"time"

	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type OutboundChannelCreated struct {
	messages.DefaultMessage
	Channel OutboundChannelCreatedChannel `json:"channel"`
}
type OutboundChannelCreatedChannel struct {
	Id          entity.ChannelId             `json:"id"`
	Name        string                       `json:"name"`
	Kind        entity.ChannelKind           `json:"kind"`
	Topic       string                       `json:"topic"`
	IsPrivate   bool                         `json:"isPrivate"`
	WorkspaceId workspace_entity.WorkspaceId `json:"workspaceId"`
	CreatedAt   time.Time                    `json:"createdAt"`
	UpdatedAt   time.Time                    `json:"updatedAt"`
	Index       int                          `json:"index"`
}

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

func (m *OutboundChannelCreated) GetActionName() messages.Action {
	return messages.OutboundChannelCreatedAction
}

func (m *OutboundChannelCreated) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}

func (m *OutboundChannelsDeleted) GetActionName() messages.Action {
	return messages.OutboundChannelsDeletedAction
}

func (m *OutboundChannelsDeleted) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}

func (m *OutboundChannelsReordered) GetActionName() messages.Action {
	return messages.OutboundChannelsReorderedAction
}

func (m *OutboundChannelsReordered) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
