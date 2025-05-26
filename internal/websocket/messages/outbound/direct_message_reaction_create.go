package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundDirectMessageReactionAdded struct {
	messages.DefaultMessage
	MessageId   string                              `json:"messageId"`
	OtherUserId string                              `json:"otherUserId"`
	Member      OutboundDirectMessageReactionMember `json:"member"`
	Reaction    string                              `json:"reaction"`
}

type OutboundDirectMessageReactionMember struct {
	UserId   string `json:"userId"`
	Username string `json:"username"`
}

func (o *OutboundDirectMessageReactionAdded) GetActionName() messages.Action {
	return messages.OutboundDirectMessageReactionAddedAction
}

func (o *OutboundDirectMessageReactionAdded) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}

type OutboundDirectMessageReactionRemoved struct {
	messages.DefaultMessage
	MessageId   string                              `json:"messageId"`
	OtherUserId string                              `json:"otherUserId"`
	Member      OutboundDirectMessageReactionMember `json:"member"`
	Reaction    string                              `json:"reaction"`
}

func (o *OutboundDirectMessageReactionRemoved) GetActionName() messages.Action {
	return messages.OutboundDirectMessageReactionRemovedAction
}

func (o *OutboundDirectMessageReactionRemoved) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
