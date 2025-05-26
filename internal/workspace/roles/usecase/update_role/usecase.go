package update_role

import (
	"context"
	"errors"

	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type UpdateRoleUseCase struct {
	roleRepository repository.RoleRepository
}

func NewUpdateRoleUseCase(roleRepository repository.RoleRepository) *UpdateRoleUseCase {
	return &UpdateRoleUseCase{roleRepository: roleRepository}
}

func (u *UpdateRoleUseCase) Execute(ctx context.Context, role entity.Role) error {
	if role.Id == "" {
		return errors.New("role ID is required")
	}

	err := u.roleRepository.Update(ctx, &role)
	if err != nil {
		return err
	}

	return nil
}
