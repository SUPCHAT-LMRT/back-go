package event

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

const (
	GroupTransferOwnershipEventType EventType = "group_transfer_ownership"
)

type GroupTransferOwnershipEvent struct {
	Group      *group_entity.Group
	NewOwnerId group_entity.GroupMemberId
}

func (e GroupTransferOwnershipEvent) Type() EventType {
	return GroupTransferOwnershipEventType
}
