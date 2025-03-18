package join_workspace_invite

import (
	"context"
	"errors"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/add_member"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/is_user_in_workspace"
	uberdig "go.uber.org/dig"
)

var (
	UserAlreadyInWorkspaceErr = errors.New("user is already in workspace")
)

type JoinWorkspaceInviteUseCaseDeps struct {
	uberdig.In
	Repository               repository.InviteLinkRepository
	IsUserInWorkspaceUseCase *is_user_in_workspace.IsUserInWorkspaceUseCase
	AddMemberUseCase         *add_member.AddMemberUseCase
}

type JoinWorkspaceInviteUseCase struct {
	deps JoinWorkspaceInviteUseCaseDeps
}

func NewJoinWorkspaceInviteUseCase(deps JoinWorkspaceInviteUseCaseDeps) *JoinWorkspaceInviteUseCase {
	return &JoinWorkspaceInviteUseCase{deps: deps}
}

func (u *JoinWorkspaceInviteUseCase) Execute(ctx context.Context, token string, user *user_entity.User) error {

	data, err := u.deps.Repository.GetInviteLinkData(ctx, token)
	if err != nil {
		return err
	}

	// Check if user is already in workspace
	isUserInWorkspace, err := u.deps.IsUserInWorkspaceUseCase.Execute(ctx, data.WorkspaceId, user.Id)
	if err != nil {
		return err
	}

	if isUserInWorkspace {
		return UserAlreadyInWorkspaceErr
	}

	err = u.deps.AddMemberUseCase.Execute(ctx, data.WorkspaceId, &entity2.WorkspaceMember{
		UserId: user.Id,
	})
	if err != nil {
		return err
	}

	err = u.deps.Repository.DeleteInviteLink(ctx, token)
	if err != nil {
		return err
	}

	return nil
}
