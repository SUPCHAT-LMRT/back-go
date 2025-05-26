package inbound

import (
	"github.com/goccy/go-json"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundDirectMessageReactionToggle struct {
	messages.DefaultMessage
	OtherUserId user_entity.UserId `json:"otherUserId"`
	MessageId   string             `json:"messageId"`
	Reaction    string             `json:"reaction"`
}

func (i *InboundDirectMessageReactionToggle) GetActionName() messages.Action {
	return messages.InboundDirectMessageReactionToggle
}

func (i *InboundDirectMessageReactionToggle) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
