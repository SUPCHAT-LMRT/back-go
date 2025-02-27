package create_workspace

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/add_member"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type CreateWorkspaceUseCase struct {
	workspaceRepository repository.WorkspaceRepository
	addMemberUseCase    *add_member.AddMemberUseCase
}

func NewCreateWorkspaceUseCase(workspaceRepository repository.WorkspaceRepository, addMemberUseCase *add_member.AddMemberUseCase) *CreateWorkspaceUseCase {
	return &CreateWorkspaceUseCase{workspaceRepository: workspaceRepository, addMemberUseCase: addMemberUseCase}
}

func (u *CreateWorkspaceUseCase) Execute(ctx context.Context, workspace *entity.Workspace, ownerMember *entity2.WorkspaceMember) error {
	err := u.workspaceRepository.Create(ctx, workspace, ownerMember)
	if err != nil {
		return err
	}

	err = u.addMemberUseCase.Execute(ctx, workspace.Id, ownerMember)
	if err != nil {
		return err
	}

	return nil
}
