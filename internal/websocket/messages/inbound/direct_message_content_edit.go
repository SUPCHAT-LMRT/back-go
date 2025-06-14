package inbound

import (
	"github.com/goccy/go-json"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundDirectMessageContentEdit struct {
	messages.DefaultMessage
	OtherUserId user_entity.UserId  `json:"otherUserId"`
	MessageId   entity.ChatDirectId `json:"messageId"`
	NewContent  string              `json:"newContent"`
}

func (i *InboundDirectMessageContentEdit) GetActionName() messages.Action {
	return messages.InboundDirectMessageContentEdit
}

func (i *InboundDirectMessageContentEdit) Encode() ([]byte, error) {
	i.DefaultMessage = messages.NewDefaultMessage(i.GetActionName())
	return json.Marshal(i)
}
