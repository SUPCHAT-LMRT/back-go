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

type DiscoverListWorkspacesUseCase struct {
	deps DiscoveryListWorkspacesUseCaseDeps
}

func NewDiscoverListWorkspacesUseCase(deps DiscoveryListWorkspacesUseCaseDeps) *DiscoverListWorkspacesUseCase {
	return &DiscoverListWorkspacesUseCase{deps: deps}
}

func (u *DiscoverListWorkspacesUseCase) Execute(ctx context.Context) ([]*DiscoverWorkspace, error) {
	publicWorkspaces, err := u.deps.WorkspaceRepository.ListPublics(ctx)
	if err != nil {
		return nil, err
	}

	discoverWorkspaces := make([]*DiscoverWorkspace, len(publicWorkspaces))
	for i, workspace := range publicWorkspaces {
		workspaceMembersCount, err := u.deps.WorkspaceMemberRepository.CountMembers(ctx, workspace.Id)
		if err != nil {
			return nil, err
		}

		discoverWorkspaces[i] = &DiscoverWorkspace{
			Id:           workspace.Id,
			Name:         workspace.Name,
			Topic:        workspace.Topic,
			OwnerId:      entity2.WorkspaceMemberId(workspace.OwnerId),
			MembersCount: workspaceMembersCount,
		}
	}

	return discoverWorkspaces, nil
}

type DiscoverWorkspace struct {
	Id           entity.WorkspaceId
	Name         string
	Topic        string
	OwnerId      entity2.WorkspaceMemberId
	MembersCount uint
}
