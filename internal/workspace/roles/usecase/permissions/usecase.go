package permissions

import (
	"context"

	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type CheckPermissionUseCase struct {
	roleRepository repository.RoleRepository
}

func NewCheckPermissionUseCase(roleRepository repository.RoleRepository) *CheckPermissionUseCase {
	return &CheckPermissionUseCase{roleRepository: roleRepository}
}

func (u *CheckPermissionUseCase) Execute(
	ctx context.Context,
	workspaceMemberId entity2.WorkspaceMemberId,
	workspaceId workspace_entity.WorkspaceId,
	permission uint64,
) (bool, error) {
	roles, err := u.roleRepository.GetRolesWithAssignmentForMember(
		ctx,
		workspaceId,
		workspaceMemberId,
	)
	if err != nil {
		return false, err
	}

	for _, role := range roles {
		if role.IsAssigned && role.HasPermission(permission) {
			return true, nil
		}
	}

	return false, nil
}
