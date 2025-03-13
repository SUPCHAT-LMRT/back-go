package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundChannelMessageReactionAdded struct {
	messages.DefaultMessage
	MessageId string                               `json:"messageId"`
	Member    OutboundChannelMessageReactionMember `json:"member"`
	Reaction  string                               `json:"reaction"`
}

type OutboundChannelMessageReactionMember struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

func (o OutboundChannelMessageReactionAdded) GetActionName() messages.Action {
	return messages.OutboundChannelMessageReactionAdded
}

func (o OutboundChannelMessageReactionAdded) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}

type OutboundChannelMessageReactionRemoved struct {
	messages.DefaultMessage
	MessageId string                               `json:"messageId"`
	Member    OutboundChannelMessageReactionMember `json:"member"`
	Reaction  string                               `json:"reaction"`
}

func (o OutboundChannelMessageReactionRemoved) GetActionName() messages.Action {
	return messages.OutboundChannelMessageReactionRemoved
}

func (o OutboundChannelMessageReactionRemoved) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
