package outbound

import (
	"github.com/goccy/go-json"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundRemovedRecentGroupChat struct {
	messages.DefaultMessage
	GroupId group_entity.GroupId `json:"groupId"`
}

func (m *OutboundRemovedRecentGroupChat) GetActionName() messages.Action {
	return messages.OutboundRecentGroupChatRemovedAction
}

func (m *OutboundRemovedRecentGroupChat) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
