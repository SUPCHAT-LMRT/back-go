package create_workspace

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/add_member"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	entity3 "github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/assign_role"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/usecase/create_role"
)

type CreateWorkspaceUseCase struct {
	workspaceRepository repository.WorkspaceRepository
	addMemberUseCase    *add_member.AddMemberUseCase
	addRoleUseCase      *create_role.CreateRoleUseCase
	assignRoleUseCase   *assign_role.AssignRoleToUserUsecase
}

func NewCreateWorkspaceUseCase(
	workspaceRepository repository.WorkspaceRepository,
	addMemberUseCase *add_member.AddMemberUseCase,
	addRoleUseCase *create_role.CreateRoleUseCase,
	assignRoleUseCase *assign_role.AssignRoleToUserUsecase,
) *CreateWorkspaceUseCase {
	return &CreateWorkspaceUseCase{
		workspaceRepository: workspaceRepository,
		addMemberUseCase:    addMemberUseCase,
		addRoleUseCase:      addRoleUseCase,
		assignRoleUseCase:   assignRoleUseCase,
	}
}

func (u *CreateWorkspaceUseCase) Execute(
	ctx context.Context,
	workspace *entity.Workspace,
	ownerMember *entity2.WorkspaceMember,
) error {
	err := u.workspaceRepository.Create(ctx, workspace, ownerMember)
	if err != nil {
		return err
	}

	err = u.addMemberUseCase.Execute(ctx, workspace.Id, ownerMember)
	if err != nil {
		return err
	}

	ownerRole := entity3.Role{
		Name:        "Owner",
		WorkspaceId: workspace.Id,
		Color:       "#f97316",
		Permissions: entity3.PermissionManageChannels |
			entity3.PermissionManageMessages |
			entity3.PermissionManageInvites |
			entity3.PermissionSendMessages |
			entity3.PermissionAttachFiles |
			entity3.PermissionPinMessages |
			entity3.PermissionMentionEveryone |
			entity3.PermissionKickMembers |
			entity3.PermissionInviteMembers |
			entity3.PermissionManageWorkspaceSettings |
			entity3.PermissionManageRoles,
	}

	err = u.addRoleUseCase.Execute(ctx, &ownerRole)
	if err != nil {
		return err
	}

	err = u.assignRoleUseCase.Execute(
		ctx,
		ownerMember.Id,
		ownerRole.Id,
		entity3.WorkspaceId(workspace.Id),
	)
	if err != nil {
		return err
	}

	return nil
}
