package get_role

import (
	"context"
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type GetRoleUseCase struct {
	roleRepository repository.RoleRepository
}

func NewGetRoleUseCase(roleRepository repository.RoleRepository) *GetRoleUseCase {
	return &GetRoleUseCase{roleRepository: roleRepository}
}

func (u *GetRoleUseCase) Execute(ctx context.Context, roleId string) (*entity.Role, error) {
	if roleId == "" {
		return nil, fmt.Errorf("roleId is required")
	}

	role, err := u.roleRepository.GetById(ctx, roleId)
	if err != nil {
		return nil, err
	}

	return role, nil
}
