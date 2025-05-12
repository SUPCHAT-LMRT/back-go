package reoder_channels

import (
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type ReorderIndexChannelsObserver interface {
	NotifyChannelReordered(channels []ChannelReorderMessage, workspaceId workspace_entity.WorkspaceId)
}

type ChannelReorderMessage struct {
	ChannelId entity.ChannelId
	NewOrder  int
}
