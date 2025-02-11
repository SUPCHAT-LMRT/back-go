package get_workspace_details

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	uberdig "go.uber.org/dig"
)

type GetWorkspaceDetailsUseCaseDeps struct {
	uberdig.In
	WorkspaceRepository repository.WorkspaceRepository
}

type GetWorkspaceDetailsUseCase struct {
	deps GetWorkspaceDetailsUseCaseDeps
}

func NewGetWorkspaceDetailsUseCase(deps GetWorkspaceDetailsUseCaseDeps) *GetWorkspaceDetailsUseCase {
	return &GetWorkspaceDetailsUseCase{deps: deps}
}

func (u *GetWorkspaceDetailsUseCase) Execute(ctx context.Context, workspaceId entity.WorkspaceId) (*WorkspaceDetails, error) {
	membersCount, err := u.deps.WorkspaceRepository.CountMembers(ctx, workspaceId)
	if err != nil {
		return nil, err
	}

	return &WorkspaceDetails{
		MembersCount: membersCount,
	}, nil
}

type WorkspaceDetails struct {
	MembersCount  uint
	ChannelsCount uint
	MessagesCount uint
}
