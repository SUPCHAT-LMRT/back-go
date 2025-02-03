package entity

import workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"

type ChannelId string

type Channel struct {
	Id          ChannelId
	Name        string
	Topic       string
	WorkspaceId workspace_entity.WorkspaceId
}

func (id ChannelId) String() string {
	return string(id)
}
