package kick_member

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type KickGroupMemberObserver interface {
	NotifyGroupMemberKicked(msg *group_entity.Group, kickedMemberId group_entity.GroupMemberId, kickedUserId user_entity.UserId)
}
