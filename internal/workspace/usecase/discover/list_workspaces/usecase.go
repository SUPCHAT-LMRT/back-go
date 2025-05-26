package list_workspaces

import (
	"context"

	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
	"github.com/supchat-lmrt/back-go/internal/user/status/usecase/get_public_status"
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
	GetPublicStatusUseCase    *get_public_status.GetPublicStatusUseCase
}

type DiscoverListWorkspacesUseCase struct {
	deps DiscoveryListWorkspacesUseCaseDeps
}

func NewDiscoverListWorkspacesUseCase(
	deps DiscoveryListWorkspacesUseCaseDeps,
) *DiscoverListWorkspacesUseCase {
	return &DiscoverListWorkspacesUseCase{deps: deps}
}

func (u *DiscoverListWorkspacesUseCase) Execute(ctx context.Context) ([]*DiscoverWorkspace, error) {
	publicWorkspaces, err := u.deps.WorkspaceRepository.ListPublics(ctx)
	if err != nil {
		return nil, err
	}

	return u.getDiscoverWorkspaces(ctx, publicWorkspaces)
}

//nolint:revive
func (u *DiscoverListWorkspacesUseCase) getDiscoverWorkspaces(
	ctx context.Context,
	publicWorkspaces []*entity.Workspace,
) ([]*DiscoverWorkspace, error) {
	discoverWorkspaces := make([]*DiscoverWorkspace, len(publicWorkspaces))
	for i, workspace := range publicWorkspaces {
		totalMembers, workspaceMembers, err := u.deps.WorkspaceMemberRepository.ListMembers(
			ctx,
			workspace.Id,
			0,
			1,
		)
		if err != nil {
			return nil, err
		}

		var onlineMembersCount uint
		for _, member := range workspaceMembers {
			memberStatus, err := u.deps.GetPublicStatusUseCase.Execute(
				ctx,
				member.UserId,
				user_status_entity.StatusOffline,
			)
			if err != nil {
				return nil, err
			}
			if memberStatus != user_status_entity.StatusOffline {
				onlineMembersCount++
			}
		}

		discoverWorkspaces[i] = &DiscoverWorkspace{
			Id:                 workspace.Id,
			Name:               workspace.Name,
			Topic:              workspace.Topic,
			OwnerId:            entity2.WorkspaceMemberId(workspace.OwnerId),
			MembersCount:       totalMembers,
			OnlineMembersCount: onlineMembersCount,
		}
	}

	return discoverWorkspaces, nil
}

type DiscoverWorkspace struct {
	Id                 entity.WorkspaceId
	Name               string
	Topic              string
	OwnerId            entity2.WorkspaceMemberId
	MembersCount       uint
	OnlineMembersCount uint
}
