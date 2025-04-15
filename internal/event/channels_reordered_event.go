package event

import (
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

const (
	ChannelsReorderedEventType EventType = "channels_reordered"
)

type ChannelsReorderedEvent struct {
	ChannelReorders []ChannelReorderMessage      `json:"channelReorders"`
	WorkspaceId     workspace_entity.WorkspaceId `json:"workspaceId"`
}

type ChannelReorderMessage struct {
	ChannelId entity.ChannelId `json:"channelId"`
	NewOrder  int              `json:"newOrder"`
}

func (e ChannelsReorderedEvent) Type() EventType {
	return ChannelsReorderedEventType
}
