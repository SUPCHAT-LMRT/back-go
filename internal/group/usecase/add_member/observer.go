package add_member

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type AddGroupMemberObserver interface {
	NotifyGroupMemberKicked(msg *group_entity.Group, inviterUserId user_entity.UserId)
}
