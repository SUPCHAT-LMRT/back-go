package repository

import (
	"context"
	"errors"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

var (
	WorkspaceNotFoundErr       = errors.New("workspace not found")
	WorkspaceMemberNotFoundErr = errors.New("workspace member not found")
	WorkspaceMemberExistsErr   = errors.New("workspace member already exists")
)

type WorkspaceRepository interface {
	Create(ctx context.Context, workspace *entity.Workspace, ownerMember *entity.WorkspaceMember) error
	GetById(ctx context.Context, id entity.WorkspaceId) (*entity.Workspace, error)
	ExistsById(ctx context.Context, id entity.WorkspaceId) (bool, error)
	List(ctx context.Context) ([]*entity.Workspace, error)
	ListByUserId(ctx context.Context, userId user_entity.UserId) ([]*entity.Workspace, error)
	ListMembers(ctx context.Context, workspaceId entity.WorkspaceId) ([]*entity.WorkspaceMember, error)
	GetMemberByUserId(ctx context.Context, workspaceId entity.WorkspaceId, userId user_entity.UserId) (*entity.WorkspaceMember, error)
	AddMember(ctx context.Context, workspaceId entity.WorkspaceId, member *entity.WorkspaceMember) error
	Update(ctx context.Context, workspace *entity.Workspace) error
	Delete(ctx context.Context, id entity.WorkspaceId) error
}
