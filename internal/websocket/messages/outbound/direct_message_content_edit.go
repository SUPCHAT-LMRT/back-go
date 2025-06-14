package outbound

import (
	"github.com/goccy/go-json"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundDirectMessageContentEdited struct {
	messages.DefaultMessage
	OtherUserId user_entity.UserId              `json:"otherUserId"`
	MessageId   chat_direct_entity.ChatDirectId `json:"messageId"`
	NewContent  string                          `json:"newContent"`
}

func (o *OutboundDirectMessageContentEdited) GetActionName() messages.Action {
	return messages.OutboundDirectMessageContentEditedAction
}

func (o *OutboundDirectMessageContentEdited) Encode() ([]byte, error) {
	o.DefaultMessage = messages.NewDefaultMessage(o.GetActionName())
	return json.Marshal(o)
}
