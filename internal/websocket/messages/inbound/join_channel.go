package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type InboundJoinChannel struct {
	messages.DefaultMessage
	ChannelId channel_entity.ChannelId `json:"channelId"`
}

func (m *InboundJoinChannel) GetActionName() messages.Action {
	return messages.InboundJoinChannelRoomAction
}

func (m *InboundJoinChannel) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
