package event

import (
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

const (
	ChannelCreatedEventType    EventType = "channels_created"
	ChannelsReorderedEventType EventType = "channels_reordered"
	ChannelsDeletedEventType   EventType = "channels_deleted"
)

type ChannelCreatedEvent struct {
	Channel *entity.Channel `json:"channel"`
}

type ChannelsReorderedEvent struct {
	ChannelReorders []ChannelReorderMessage      `json:"channelReorders"`
	WorkspaceId     workspace_entity.WorkspaceId `json:"workspaceId"`
}

type ChannelsDeletedEvent struct {
	ChannelId   entity.ChannelId             `json:"channelId"`
	WorkspaceId workspace_entity.WorkspaceId `json:"workspaceId"`
}

type ChannelReorderMessage struct {
	ChannelId entity.ChannelId `json:"channelId"`
	NewOrder  int              `json:"newOrder"`
}

func (e ChannelsReorderedEvent) Type() EventType {
	return ChannelsReorderedEventType
}

func (c ChannelsDeletedEvent) Type() EventType {
	return ChannelsDeletedEventType
}

func (c ChannelCreatedEvent) Type() EventType {
	return ChannelCreatedEventType
}
