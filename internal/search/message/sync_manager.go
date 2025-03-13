package message

import (
	"context"
)

type SearchMessageSyncManager interface {
	AddMessage(ctx context.Context, message *SearchMessage) error
	SyncLoop(ctx context.Context)
}
