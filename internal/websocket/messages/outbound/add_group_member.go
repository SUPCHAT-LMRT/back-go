package outbound

import (
	"github.com/goccy/go-json"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
)

type OutboundAddGroupMemberChat struct {
	messages.DefaultMessage
	GroupId group_entity.GroupId              `json:"groupId"`
	Member  *OutboundAddGroupMemberChatMember `json:"member"`
}

type OutboundAddGroupMemberChatMember struct {
	Id           group_entity.GroupMemberId `json:"id"`
	UserId       user_entity.UserId         `json:"userId"`
	UserName     string                     `json:"userName"`
	IsGroupOwner bool                       `json:"isGroupOwner"`
	Status       entity.Status              `json:"status"`
}

func (m *OutboundAddGroupMemberChat) GetActionName() messages.Action {
	return messages.OutboundGroupMemberAddedAction
}

func (m *OutboundAddGroupMemberChat) Encode() ([]byte, error) {
	m.DefaultMessage = messages.NewDefaultMessage(m.GetActionName())
	return json.Marshal(m)
}
