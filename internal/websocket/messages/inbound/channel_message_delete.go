package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type InboundChannelMessageDelete struct {
	messages.DefaultMessage
	ChannelId channel_entity.ChannelId `json:"channelId"`
	MessageId entity.ChannelMessageId  `json:"messageId"`
}

func (i *InboundChannelMessageDelete) GetActionName() messages.Action {
	return messages.InboundChannelMessageDelete
}

func (i *InboundChannelMessageDelete) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
