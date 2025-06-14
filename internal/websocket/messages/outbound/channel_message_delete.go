package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type OutboundChannelMessageDeleted struct {
	messages.DefaultMessage
	ChannelId channel_entity.ChannelId                `json:"channelId"`
	MessageId channel_message_entity.ChannelMessageId `json:"messageId"`
}

func (o *OutboundChannelMessageDeleted) GetActionName() messages.Action {
	return messages.OutboundChannelMessageDeletedAction
}

func (o *OutboundChannelMessageDeleted) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
