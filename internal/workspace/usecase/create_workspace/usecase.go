package create_workspace

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type CreateWorkspaceUseCase struct {
	workspaceRepository repository.WorkspaceRepository
}

func NewCreateWorkspaceUseCase(workspaceRepository repository.WorkspaceRepository) *CreateWorkspaceUseCase {
	return &CreateWorkspaceUseCase{workspaceRepository: workspaceRepository}
}

func (u *CreateWorkspaceUseCase) Execute(ctx context.Context, workspace *entity.Workspace, ownerMember *entity.WorkspaceMember) error {
	return u.workspaceRepository.Create(ctx, workspace, ownerMember)
}
