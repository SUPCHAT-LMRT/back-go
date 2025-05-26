package strategies

import (
	"context"

	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

type DefaultGroupNameStrategy interface {
	Handle(
		ctx context.Context,
		group *group_entity.Group,
		members []*group_entity.GroupMember,
	) (string, error)
}
