package inbound

import (
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type InboundJoinDirectRoom struct {
	messages.DefaultMessage
	ChannelId channel_entity.ChannelId `json:"channelId"`
}

func (m InboundJoinDirectRoom) GetActionName() messages.Action {
	return messages.InboundJoinChannelRoomAction
}
