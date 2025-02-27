package list_workspaces

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	repository2 "github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type DiscoveryListWorkspacesUseCase struct {
	workspaceRepository       repository.WorkspaceRepository
	workspaceMemberRepository repository2.WorkspaceMemberRepository
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
		workspaceMembersCount, err := u.workspaceMemberRepository.CountMembers(ctx, workspace.Id)
		if err != nil {
			return nil, err
		}

		discoverWorkspaces[i] = &DiscoveryWorkspace{
			Id:           workspace.Id,
			Name:         workspace.Name,
			OwnerId:      entity2.WorkspaceMemberId(workspace.OwnerId),
			MembersCount: workspaceMembersCount,
		}
	}

	return discoverWorkspaces, nil
}

type DiscoveryWorkspace struct {
	Id           entity.WorkspaceId
	Name         string
	OwnerId      entity2.WorkspaceMemberId
	MembersCount uint
}
