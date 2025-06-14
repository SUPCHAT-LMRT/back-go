package outbound

import (
	"github.com/goccy/go-json"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundAddRecentDirectChat struct {
	messages.DefaultMessage
	OtherUserId user_entity.UserId `json:"otherUserId"`
	ChatName    string             `json:"chatName"`
}

func (m *OutboundAddRecentDirectChat) GetActionName() messages.Action {
	return messages.OutboundRecentDirectChatAddedAction
}

func (m *OutboundAddRecentDirectChat) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}

type OutboundAddRecentGroupChat struct {
	messages.DefaultMessage
	GroupId  group_entity.GroupId `json:"groupId"`
	ChatName string               `json:"chatName"`
}

func (m *OutboundAddRecentGroupChat) GetActionName() messages.Action {
	return messages.OutboundRecentGroupChatAddedAction
}

func (m *OutboundAddRecentGroupChat) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
