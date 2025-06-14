package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type OutboundChannelMessageContentEdited struct {
	messages.DefaultMessage
	ChannelId  channel_entity.ChannelId                `json:"channelId"`
	MessageId  channel_message_entity.ChannelMessageId `json:"messageId"`
	NewContent string                                  `json:"newContent"`
}

func (o *OutboundChannelMessageContentEdited) GetActionName() messages.Action {
	return messages.OutboundChannelMessageContentEditedAction
}

func (o *OutboundChannelMessageContentEdited) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
