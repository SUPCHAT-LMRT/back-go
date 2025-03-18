package outbound

import (
	"github.com/goccy/go-json"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundAddRecentDirectChat struct {
	messages.DefaultMessage
	OtherUserId user_entity.UserId `json:"otherUserId"`
	ChatName    string             `json:"chatName"`
}

func (m OutboundAddRecentDirectChat) GetActionName() messages.Action {
	return messages.OutboundRecentDirectChatAddedAction
}

func (m OutboundAddRecentDirectChat) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
