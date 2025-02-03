package list_workspaces

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type ListWorkspacesUseCase struct {
	workspaceRepository repository.WorkspaceRepository
}

func NewListWorkspacesUseCase(workspaceRepository repository.WorkspaceRepository) *ListWorkspacesUseCase {
	return &ListWorkspacesUseCase{workspaceRepository: workspaceRepository}
}

func (u *ListWorkspacesUseCase) Execute(ctx context.Context, userId user_entity.UserId) ([]*entity.Workspace, error) {
	return u.workspaceRepository.ListByUserId(ctx, userId)
}
