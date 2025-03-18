package repository

import (
	"context"
	"errors"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

var (
	WorkspaceMemberNotFoundErr = errors.New("workspace member not found")
	WorkspaceMemberExistsErr   = errors.New("workspace member already exists")
)

type WorkspaceMemberRepository interface {
	ListMembers(ctx context.Context, workspaceId entity.WorkspaceId, limit, page int) (totalMembers uint, members []*entity2.WorkspaceMember, err error)
	CountMembers(ctx context.Context, workspaceId entity.WorkspaceId) (uint, error)
	GetMemberByUserId(ctx context.Context, workspaceId entity.WorkspaceId, userId user_entity.UserId) (*entity2.WorkspaceMember, error)
	IsMemberExists(ctx context.Context, workspaceId entity.WorkspaceId, userId user_entity.UserId) (bool, error)
	AddMember(ctx context.Context, workspaceId entity.WorkspaceId, member *entity2.WorkspaceMember) error
}
