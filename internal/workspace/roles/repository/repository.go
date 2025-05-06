package repository

import (
	"context"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
)

type RoleRepository interface {
	Create(ctx context.Context, role *entity.Role) error
	GetById(ctx context.Context, roleId string) (*entity.Role, error)
	GetList(ctx context.Context, workspaceId string) ([]*entity.Role, error)
	Update(ctx context.Context, role *entity.Role) error
	Delete(ctx context.Context, roleId string) error
	AssignRoleToUser(ctx context.Context, userId string, roleId string, workspaceId string) error
	DessassignRoleFromUser(ctx context.Context, userId string, roleId string, workspaceId string) error
	GetRolesWithAssignmentForMember(ctx context.Context, workspaceId workspace_entity.WorkspaceId, workspaceMemberId entity2.WorkspaceMemberId) ([]*entity.Role, error)
}
