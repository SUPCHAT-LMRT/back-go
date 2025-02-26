package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundMessageReactionAdded struct {
	messages.DefaultMessage
	MessageId string                        `json:"messageId"`
	Member    OutboundMessageReactionMember `json:"member"`
	Reaction  string                        `json:"reaction"`
}

type OutboundMessageReactionMember struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

func (o OutboundMessageReactionAdded) GetActionName() messages.Action {
	return messages.OutboundChannelMessageReactionAdded
}

func (o OutboundMessageReactionAdded) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}

type OutboundMessageReactionRemoved struct {
	messages.DefaultMessage
	MessageId string                        `json:"messageId"`
	Member    OutboundMessageReactionMember `json:"member"`
	Reaction  string                        `json:"reaction"`
}

func (o OutboundMessageReactionRemoved) GetActionName() messages.Action {
	return messages.OutboundChannelMessageReactionRemoved
}

func (o OutboundMessageReactionRemoved) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
