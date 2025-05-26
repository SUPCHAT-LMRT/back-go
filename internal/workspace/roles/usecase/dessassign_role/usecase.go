package dessassign_role

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type DessassignRoleFromUserUsecase struct {
	roleRepository repository.RoleRepository
}

func NewDessassignRoleFromUserUsecase(
	roleRepository repository.RoleRepository,
) *DessassignRoleFromUserUsecase {
	return &DessassignRoleFromUserUsecase{roleRepository: roleRepository}
}

func (u *DessassignRoleFromUserUsecase) Execute(
	ctx context.Context,
	userId string,
	roleId string,
	workspaceId string,
) error {
	err := u.roleRepository.DessassignRoleFromUser(ctx, userId, roleId, workspaceId)
	if err != nil {
		return err
	}

	return nil
}
