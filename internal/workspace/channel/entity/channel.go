package entity

import (
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"time"
)

type ChannelId string

type ChannelKind string

const (
	ChannelKindUnknown ChannelKind = ""
	ChannelKindText    ChannelKind = "text"
	ChannelKindVoice   ChannelKind = "voice"
)

type Channel struct {
	Id          ChannelId
	Name        string
	Topic       string
	Kind        ChannelKind
	WorkspaceId workspace_entity.WorkspaceId
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (id ChannelId) String() string {
	return string(id)
}

func (k ChannelKind) String() string {
	return string(k)
}
