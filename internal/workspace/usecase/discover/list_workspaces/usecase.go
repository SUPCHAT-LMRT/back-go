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

func (u *DiscoveryListWorkspacesUseCase) Execute(ctx context.Context) ([]*DiscoveryWorkspace, error) {
	publicWorkspaces, err := u.workspaceRepository.ListPublics(ctx)
	if err != nil {
		return nil, err
	}

	discoverWorkspaces := make([]*DiscoveryWorkspace, len(publicWorkspaces))
	for i, workspace := range publicWorkspaces {
		workspaceMembersCount, err := u.workspaceRepository.CountMembers(ctx, workspace.Id)
		if err != nil {
			return nil, err
		}

		discoverWorkspaces[i] = &DiscoveryWorkspace{
			Id:           workspace.Id,
			Name:         workspace.Name,
			OwnerId:      entity.WorkspaceMemberId(workspace.OwnerId),
			MembersCount: workspaceMembersCount,
		}
	}

	return discoverWorkspaces, nil
}

type DiscoveryWorkspace struct {
	Id           entity.WorkspaceId
	Name         string
	OwnerId      entity.WorkspaceMemberId
	MembersCount uint
}
