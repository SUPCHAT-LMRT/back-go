package list_mentionnable_user

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/list_user_private_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/get_user_by_workspace_member_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/list_workspace_members"
	repository2 "github.com/supchat-lmrt/back-go/internal/workspace/repository"
)

type (
	MentionnableUser struct {
		Id       entity.UserId
		Username string
	}

	ListMentionnableUserUseCase struct {
		channelRepo                       repository.ChannelRepository
		workspaceRepo                     repository2.WorkspaceRepository
		listPrivateChannelMembersUsecase  *list_user_private_channel.ListPrivateChannelMembersUseCase
		getUserByWorkspaceMemberIdUseCase *get_user_by_workspace_member_id.GetUserByWorkspaceMemberIdUseCase
		listWorkspaceMembersUseCase       *list_workspace_members.ListWorkspaceMembersUseCase
	}
)

func NewListMentionnableUserUseCase(
	channelRepo repository.ChannelRepository,
	workspaceRepo repository2.WorkspaceRepository,
	listPrivateChannelMembersUsecase *list_user_private_channel.ListPrivateChannelMembersUseCase,
	getUserByWorkspaceMemberIdUseCase *get_user_by_workspace_member_id.GetUserByWorkspaceMemberIdUseCase,
	listWorkspaceMembersUseCase *list_workspace_members.ListWorkspaceMembersUseCase,
) *ListMentionnableUserUseCase {
	return &ListMentionnableUserUseCase{
		channelRepo:                       channelRepo,
		workspaceRepo:                     workspaceRepo,
		listPrivateChannelMembersUsecase:  listPrivateChannelMembersUsecase,
		getUserByWorkspaceMemberIdUseCase: getUserByWorkspaceMemberIdUseCase,
		listWorkspaceMembersUseCase:       listWorkspaceMembersUseCase,
	}
}

func (u *ListMentionnableUserUseCase) Execute(ctx context.Context, channelId channel_entity.ChannelId) ([]MentionnableUser, error) {
	channel, err := u.channelRepo.GetById(ctx, channelId)
	if err != nil {
		return nil, err
	}

	if channel.IsPrivate {
		memberIds, err := u.listPrivateChannelMembersUsecase.Execute(ctx, channelId)
		if err != nil {
			return nil, err
		}
		var users []MentionnableUser
		for _, memberId := range memberIds {
			user, err := u.getUserByWorkspaceMemberIdUseCase.Execute(ctx, memberId)
			if err != nil {
				return nil, err
			}
			users = append(users, MentionnableUser{
				Id:       user.Id,
				Username: user.FullName(),
			})
		}
		return users, nil
	}

	_, workspaceMembers, err := u.listWorkspaceMembersUseCase.Execute(ctx, channel.WorkspaceId, 1000, 1)
	if err != nil {
		return nil, err
	}
	var users []MentionnableUser
	for _, member := range workspaceMembers {
		user, err := u.getUserByWorkspaceMemberIdUseCase.Execute(ctx, member.Id)
		if err != nil {
			return nil, err
		}
		users = append(users, MentionnableUser{
			Id:       user.Id,
			Username: user.FullName(),
		})
	}
	return users, nil
}
