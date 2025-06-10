package get_list_roles

import (
	"context"
	"errors"

	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type GetListRolesUseCase struct {
	roleRepository repository.RoleRepository
}

func NewGetListRolesUseCase(roleRepository repository.RoleRepository) *GetListRolesUseCase {
	return &GetListRolesUseCase{roleRepository: roleRepository}
}

func (u *GetListRolesUseCase) Execute(
	ctx context.Context,
	workspaceId string,
) ([]*entity.Role, error) {
	if workspaceId == "" {
		return nil, errors.New("workspaceId is required")
	}

	roles, err := u.roleRepository.GetList(ctx, workspaceId)
	if err != nil {
		return nil, err
	}

	return roles, nil
}
