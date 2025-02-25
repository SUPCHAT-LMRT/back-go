package websocket

import (
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type JoinChannel struct {
	DefaultMessage
	ChannelId channel_entity.ChannelId `json:"channelId"`
}

func (m JoinChannel) GetActionName() string {
	return InboundJoinChannelRoomAction
}
