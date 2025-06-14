package transfer_ownership

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

type NotifyTransferGroupOwnershipObserver struct {
	eventBus *event.EventBus
}

func NewNotifyTransferGroupOwnershipObserver(eventBus *event.EventBus) TransferGroupOwnershipObserver {
	return &NotifyTransferGroupOwnershipObserver{eventBus: eventBus}
}

func (o NotifyTransferGroupOwnershipObserver) NotifyOwnershipTransferred(group *group_entity.Group, newOwnerId group_entity.GroupMemberId) {
	o.eventBus.Publish(&event.GroupTransferOwnershipEvent{
		Group:      group,
		NewOwnerId: newOwnerId,
	})
}
