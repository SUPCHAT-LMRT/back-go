package transfer_ownership

import group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"

type TransferGroupOwnershipObserver interface {
	NotifyOwnershipTransferred(group *group_entity.Group, newOwnerId group_entity.GroupMemberId)
}
