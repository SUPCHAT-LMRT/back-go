package assign_role

import (
	"context"

	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type AssignRoleToUserUsecase struct {
	roleRepository repository.RoleRepository
}

func NewAssignRoleToUserUsecase(roleRepository repository.RoleRepository) *AssignRoleToUserUsecase {
	return &AssignRoleToUserUsecase{roleRepository: roleRepository}
}

func (u *AssignRoleToUserUsecase) Execute(
	ctx context.Context,
	workspaceMemberId entity2.WorkspaceMemberId,
	roleId entity.RoleId,
	workspaceId entity.WorkspaceId,
) error {
	err := u.roleRepository.AssignRoleToUser(
		ctx,
		workspaceMemberId,
		roleId,
		workspace_entity.WorkspaceId(workspaceId),
	)
	if err != nil {
		return err
	}

	return nil
}
