package group

import (
	"context"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

type SearchGroupSyncManager interface {
	CreateIndexIfNotExists(context.Context) error
	AddGroup(ctx context.Context, group *SearchGroup) error
	RemoveGroup(ctx context.Context, groupId group_entity.GroupId) error
	Sync(ctx context.Context)
	SyncLoop(ctx context.Context)
}
