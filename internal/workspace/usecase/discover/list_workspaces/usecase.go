package list_workspaces

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type DiscoveryListWorkspacesUseCase struct {
	workspaceRepository repository.WorkspaceRepository
}

func NewDiscoveryListWorkspacesUseCase(workspaceRepository repository.WorkspaceRepository) *DiscoveryListWorkspacesUseCase {
	return &DiscoveryListWorkspacesUseCase{workspaceRepository: workspaceRepository}
}

func (u *DiscoveryListWorkspacesUseCase) Execute(ctx context.Context) ([]*entity.Workspace, error) {
	return u.workspaceRepository.ListPublics(ctx)
}
