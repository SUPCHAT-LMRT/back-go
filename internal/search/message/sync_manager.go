package message

import (
	"context"
)

type SearchMessageSyncManager interface {
	CreateIndexIfNotExists(ctx context.Context) error
	AddMessage(ctx context.Context, message *SearchMessage) error
	Sync(ctx context.Context)
	SyncLoop(ctx context.Context)
}
