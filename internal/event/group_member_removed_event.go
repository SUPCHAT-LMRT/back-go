package event

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

const (
	GroupMemberRemovedEventType EventType = "group_member_removed"
)

type GroupMemberRemovedEvent struct {
	Group           *group_entity.Group
	RemovedMemberId group_entity.GroupMemberId
	RemovedUserId   user_entity.UserId
}

func (e GroupMemberRemovedEvent) Type() EventType {
	return GroupMemberRemovedEventType
}
