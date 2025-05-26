package channel

import (
	"time"

	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type SearchChannelKind string

const (
	SearchChannelKindUnknown SearchChannelKind = ""
	SearchChannelKindVoice   SearchChannelKind = "voice"
	SearchChannelKindText    SearchChannelKind = "text"
)

type SearchChannel struct {
	Id          channel_entity.ChannelId
	Name        string
	Topic       string
	Kind        SearchChannelKind
	IsPrivate   bool
	Members     []string
	WorkspaceId workspace_entity.WorkspaceId
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
