package create_role

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type CreateRoleUseCase struct {
	roleRepository repository.RoleRepository
}

func NewCreateRoleUseCase(roleRepository repository.RoleRepository) *CreateRoleUseCase {
	return &CreateRoleUseCase{roleRepository: roleRepository}
}

func (u *CreateRoleUseCase) Execute(ctx context.Context, role *entity.Role) error {
	id, err := u.roleRepository.Create(ctx, role)
	if err != nil {
		return err
	}

	role.Id = id
	return nil
}
