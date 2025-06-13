package create_group

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

type GroupCreatedObserver interface {
	NotifyGroupMemberAdded(msg *group_entity.Group)
}
