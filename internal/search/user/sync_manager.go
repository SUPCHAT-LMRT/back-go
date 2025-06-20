package user

import (
	"context"
)

type SearchUserSyncManager interface {
	CreateIndexIfNotExists(context.Context) error
	AddUser(ctx context.Context, channel *SearchUser) error
	UpdateUser(ctx context.Context, channel *SearchUser) error
	RemoveUser(ctx context.Context, channel *SearchUser) error
	Sync(ctx context.Context)
	SyncLoop(ctx context.Context)
}
