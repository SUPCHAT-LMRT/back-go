package inbound

import (
	"github.com/goccy/go-json"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type InboundJoinGroupRoom struct {
	messages.DefaultMessage
	GroupId group_entity.GroupId `json:"groupId"`
}

func (m *InboundJoinGroupRoom) GetActionName() messages.Action {
	return messages.InboundJoinGroupRoomAction
}

func (m *InboundJoinGroupRoom) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
