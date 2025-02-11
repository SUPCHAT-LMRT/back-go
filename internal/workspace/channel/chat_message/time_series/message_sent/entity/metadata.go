package entity

import (
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"time"
)

type MessageSentMetadata struct {
	WorkspaceId    workspace_entity.WorkspaceId
	ChannelId      channel_entity.ChannelId
	AuthorMemberId workspace_entity.WorkspaceMemberId
}

type MessageSent struct {
	SentAt   time.Time
	Count    uint
	Metadata MessageSentMetadata
}
