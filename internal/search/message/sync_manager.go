package message

import (
	"context"
)

type SearchMessageSyncManager interface {
	CreateIndexIfNotExists(ctx context.Context) error
	AddMessage(ctx context.Context, message *SearchMessage) error
	SyncLoop(ctx context.Context)
}
