package list_workspaces

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	repository2 "github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	uberdig "go.uber.org/dig"
)

type DiscoveryListWorkspacesUseCaseDeps struct {
	uberdig.In
	WorkspaceRepository       repository.WorkspaceRepository
	WorkspaceMemberRepository repository2.WorkspaceMemberRepository
}

type DiscoveryListWorkspacesUseCase struct {
	deps DiscoveryListWorkspacesUseCaseDeps
}

func NewDiscoveryListWorkspacesUseCase(deps DiscoveryListWorkspacesUseCaseDeps) *DiscoveryListWorkspacesUseCase {
	return &DiscoveryListWorkspacesUseCase{deps: deps}
}

func (u *DiscoveryListWorkspacesUseCase) Execute(ctx context.Context) ([]*DiscoveryWorkspace, error) {
	publicWorkspaces, err := u.deps.WorkspaceRepository.ListPublics(ctx)
	if err != nil {
		return nil, err
	}

	discoverWorkspaces := make([]*DiscoveryWorkspace, len(publicWorkspaces))
	for i, workspace := range publicWorkspaces {
		workspaceMembersCount, err := u.deps.WorkspaceMemberRepository.CountMembers(ctx, workspace.Id)
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
