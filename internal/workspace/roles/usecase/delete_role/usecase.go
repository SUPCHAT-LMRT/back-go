package delete_role

import (
	"context"
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type DeleteRoleUseCase struct {
	roleRepository repository.RoleRepository
}

func NewDeleteRoleUseCase(roleRepository repository.RoleRepository) *DeleteRoleUseCase {
	return &DeleteRoleUseCase{roleRepository: roleRepository}
}

func (u *DeleteRoleUseCase) Execute(ctx context.Context, roleId string) error {
	if roleId == "" {
		return fmt.Errorf("roleId is required")
	}

	err := u.roleRepository.Delete(ctx, roleId)
	if err != nil {
		return err
	}

	return nil
}
