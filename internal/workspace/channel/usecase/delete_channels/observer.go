package delete_channels

import (
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type DeleteSpecifyChannelsObserver interface {
	NotifyChannelsDeleted(channelId entity.ChannelId, workspace workspace_entity.WorkspaceId)
}
