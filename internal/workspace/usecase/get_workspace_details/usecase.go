package get_workspace_details

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/count_messages_by_workspace"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/count_channels"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	repository2 "github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/repository"
	uberdig "go.uber.org/dig"
)

type GetWorkspaceDetailsUseCaseDeps struct {
	uberdig.In
	WorkspaceRepository       repository.WorkspaceRepository
	WorskpaceMemberRepository repository2.WorkspaceMemberRepository
	CountChannelsUseCase            *count_channels.CountChannelsUseCase
	CountMessagesByWorkspaceUseCase *count_messages_by_workspace.CountMessagesByWorkspaceUseCase
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
	membersCount, err := u.deps.WorskpaceMemberRepository.CountMembers(ctx, workspaceId)
	if err != nil {
		return nil, err
	}

	channelsCount, err := u.deps.CountChannelsUseCase.Execute(ctx, workspaceId)
	if err != nil {
		return nil, err
	}

	messagesCount, err := u.deps.CountMessagesByWorkspaceUseCase.Execute(ctx, workspaceId)
	if err != nil {
		return nil, err
	}

	return &WorkspaceDetails{
		Id:            workspace.Id,
		Name:          workspace.Name,
		Type:          workspace.Type,
		MembersCount:  membersCount,
		ChannelsCount: channelsCount,
		MessagesCount: messagesCount,
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
