package entity

import "github.com/supchat-lmrt/back-go/internal/user/entity"

type (
	WorkspaceType string
	WorkspaceId   string
)

const (
	WorkspaceTypePrivate WorkspaceType = "PRIVATE"
	WorkspaceTypePublic  WorkspaceType = "PUBLIC"
)

type Workspace struct {
	Id      WorkspaceId
	Name    string
	Topic   string
	Type    WorkspaceType
	OwnerId entity.UserId
	TrucID  entity.UserId
}

func (id WorkspaceId) String() string {
	return string(id)
}
