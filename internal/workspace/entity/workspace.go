package entity

import "github.com/supchat-lmrt/back-go/internal/user/entity"

type (
	WorkspaceType string
	WorkspaceId   string
)

const (
	WorkspaceTypePrivate WorkspaceType = "private"
	WorkspaceTypePublic  WorkspaceType = "public"
)

type Workspace struct {
	Id      WorkspaceId
	Name    string
	Type    WorkspaceType
	OwnerId entity.UserId
}

func (id WorkspaceId) String() string {
	return string(id)
}
