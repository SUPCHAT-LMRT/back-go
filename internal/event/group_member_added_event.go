package event

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

const (
	GroupMemberAddedEventType EventType = "group_member_added"
)

type GroupMemberAddedEvent struct {
	Message       *group_entity.Group
	InvitedUserId user_entity.UserId
}

func (e GroupMemberAddedEvent) Type() EventType {
	return GroupMemberAddedEventType
}
