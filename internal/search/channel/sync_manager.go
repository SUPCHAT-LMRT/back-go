package channel

import (
	"context"
)

type SearchChannelSyncManager interface {
	CreateIndexIfNotExists(context.Context) error
	AddChannel(ctx context.Context, channel *SearchChannel) error
	RemoveChannel(ctx context.Context, channel *SearchChannel) error
	SyncLoop(ctx context.Context)
}
