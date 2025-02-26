package outbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundMessageReactionCreate struct {
	messages.DefaultMessage
	ReactionId string                              `json:"reactionId"`
	MessageId  string                              `json:"messageId"`
	Member     OutboundMessageReactionCreateMember `json:"member"`
	Reaction   string                              `json:"reaction"`
}

type OutboundMessageReactionCreateMember struct {
	UserId        string `json:"userId"`
	Username      string `json:"username"`
	WorkspaceName string `json:"workspaceName"`
}

func (o OutboundMessageReactionCreate) GetActionName() messages.Action {
	return messages.OutboundChannelMessageReactionCreate
}

func (o OutboundMessageReactionCreate) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
