package entity

import (
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type WorkspaceMemberId string

type WorkspaceMember struct {
	Id          WorkspaceMemberId
	WorkspaceId entity2.WorkspaceId
	UserId      entity.UserId
}

func (id WorkspaceMemberId) String() string {
	return string(id)
}
