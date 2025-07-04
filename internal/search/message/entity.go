package message

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"time"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type SearchMessageKind string

const (
	SearchMessageKindChannelMessage SearchMessageKind = "channel"
	SearchMessageKindDirectMessage  SearchMessageKind = "direct"
	SearchMessageGroupMessage       SearchMessageKind = "group"
)

type SearchMessage struct {
	Id       string
	Content  string
	AuthorId user_entity.UserId
	Kind     SearchMessageKind
	// Data is depending on the kind (SearchMessageChannelData, SearchMessageDirectData, SearchMessageGroupData)
	Data      any
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SearchMessageChannelData struct {
	ChannelId   channel_entity.ChannelId
	WorkspaceId workspace_entity.WorkspaceId
}

type SearchMessageDirectData struct {
	User1 user_entity.UserId
	User2 user_entity.UserId
}

type SearchMessageGroupData struct {
	GroupId group_entity.GroupId
}
