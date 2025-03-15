package channel

import (
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"time"
)

type SearchChannelKind string

const (
	SearchChannelKindVoiceMessage SearchChannelKind = "voice"
	SearchChannelKindTextMessage  SearchChannelKind = "text"
)

type SearchChannel struct {
	Id          channel_entity.ChannelId
	Name        string
	Topic       string
	Kind        SearchChannelKind
	WorkspaceId workspace_entity.WorkspaceId
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
