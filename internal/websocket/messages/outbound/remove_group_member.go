package outbound

import (
	"github.com/goccy/go-json"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundRemoveGroupMemberChat struct {
	messages.DefaultMessage
	GroupId  group_entity.GroupId       `json:"groupId"`
	MemberId group_entity.GroupMemberId `json:"memberId"`
	UserId   user_entity.UserId         `json:"userId"`
}

func (m *OutboundRemoveGroupMemberChat) GetActionName() messages.Action {
	return messages.OutboundGroupMemberRemovedAction
}

func (m *OutboundRemoveGroupMemberChat) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
