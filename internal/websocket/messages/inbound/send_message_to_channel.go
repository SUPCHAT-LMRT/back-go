package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type InboundSendMessageToChannel struct {
	messages.DefaultMessage
	Content   string                   `json:"content"`
	ChannelId channel_entity.ChannelId `json:"channelId"`
}

func (m *InboundSendMessageToChannel) GetActionName() messages.Action {
	return messages.InboundSendChannelMessageAction
}

func (m *InboundSendMessageToChannel) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
