package assign_role

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type AssignRoleToUserUsecase struct {
	roleRepository repository.RoleRepository
}

func NewAssignRoleToUserUsecase(roleRepository repository.RoleRepository) *AssignRoleToUserUsecase {
	return &AssignRoleToUserUsecase{roleRepository: roleRepository}
}
func (u *AssignRoleToUserUsecase) Execute(ctx context.Context, userId string, roleId string, workspaceId string) error {
	err := u.roleRepository.AssignRoleToUser(ctx, userId, roleId, workspaceId)
	if err != nil {
		return err
	}

	return nil
}
