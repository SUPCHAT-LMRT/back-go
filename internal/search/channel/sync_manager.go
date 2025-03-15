package channel

import (
	"context"
)

type SearchChannelSyncManager interface {
	CreateIndexIfNotExists(context.Context) error
	AddChannel(ctx context.Context, channel *SearchChannel) error
	RemoveChannel(ctx context.Context, channel *SearchChannel) error
	Sync(ctx context.Context)
	SyncLoop(ctx context.Context)
}
