package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type InboundChannelMessageContentEdit struct {
	messages.DefaultMessage
	ChannelId  channel_entity.ChannelId `json:"channelId"`
	MessageId  entity.ChannelMessageId  `json:"messageId"`
	NewContent string                   `json:"newContent"`
}

func (i *InboundChannelMessageContentEdit) GetActionName() messages.Action {
	return messages.InboundChannelMessageContentEdit
}

func (i *InboundChannelMessageContentEdit) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
