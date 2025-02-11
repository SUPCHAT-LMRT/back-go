package entity

import "github.com/supchat-lmrt/back-go/internal/user/entity"

type WorkspaceMemberId string

type WorkspaceMember struct {
	Id          WorkspaceMemberId
	WorkspaceId WorkspaceId
	UserId      entity.UserId
	Pseudo      string
}

func (id WorkspaceMemberId) String() string {
	return string(id)
}
