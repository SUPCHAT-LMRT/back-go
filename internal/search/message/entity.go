package message

import (
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
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
	ChannelId string
}

type SearchMessageDirectData struct {
	OtherUserId user_entity.UserId
}

type SearchMessageGroupData struct {
	GroupId string
}
