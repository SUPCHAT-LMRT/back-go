package channel

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
)

type SearchChannelSyncManager interface {
	CreateIndexIfNotExists(context.Context) error
	AddChannel(ctx context.Context, channel *SearchChannel) error
	RemoveChannel(ctx context.Context, channelId entity.ChannelId) error
	Sync(ctx context.Context)
	SyncLoop(ctx context.Context)
}
