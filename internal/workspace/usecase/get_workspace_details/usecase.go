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
	workspace, err := u.deps.WorkspaceRepository.GetById(ctx, workspaceId)
	if err != nil {
		return nil, err
	}

	membersCount, err := u.deps.WorkspaceRepository.CountMembers(ctx, workspaceId)
	if err != nil {
		return nil, err
	}

	return &WorkspaceDetails{
		Id:           workspace.Id,
		Name:         workspace.Name,
		Type:         workspace.Type,
		MembersCount: membersCount,
	}, nil
}

type WorkspaceDetails struct {
	Id            entity.WorkspaceId
	Name          string
	Type          entity.WorkspaceType
	MembersCount  uint
	ChannelsCount uint
	MessagesCount uint
}
