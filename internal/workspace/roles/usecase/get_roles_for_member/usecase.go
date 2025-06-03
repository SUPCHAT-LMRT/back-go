package get_roles_for_member

import (
	"context"

	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/repository"
)

type GetRolesForMemberUsecase struct {
	roleRepository repository.RoleRepository
}

func NewGetRolesForMemberUsecase(
	roleRepository repository.RoleRepository,
) *GetRolesForMemberUsecase {
	return &GetRolesForMemberUsecase{roleRepository: roleRepository}
}

func (u *GetRolesForMemberUsecase) Execute(
	ctx context.Context,
	workspaceId workspace_entity.WorkspaceId,
	workspaceMemberId entity2.WorkspaceMemberId,
) ([]*entity.Role, error) {
	return u.roleRepository.GetRolesWithAssignmentForMember(ctx, workspaceId, workspaceMemberId)
}
