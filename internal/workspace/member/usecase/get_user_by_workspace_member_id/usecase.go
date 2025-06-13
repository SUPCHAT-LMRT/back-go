package get_user_by_workspace_member_id

import (
	"context"

	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
	uberdig "go.uber.org/dig"
)

type GetUserByWorkspaceMemberIdUseCaseDeps struct {
	uberdig.In
	WorkspaceMemberRepository repository.WorkspaceMemberRepository
	GetUserByIdUseCase        *get_by_id.GetUserByIdUseCase
}

type GetUserByWorkspaceMemberIdUseCase struct {
	deps GetUserByWorkspaceMemberIdUseCaseDeps
}

func NewGetUserByWorkspaceMemberIdUseCase(
	deps GetUserByWorkspaceMemberIdUseCaseDeps,
) *GetUserByWorkspaceMemberIdUseCase {
	return &GetUserByWorkspaceMemberIdUseCase{deps: deps}
}

func (u *GetUserByWorkspaceMemberIdUseCase) Execute(
	ctx context.Context,
	memberId entity.WorkspaceMemberId,
) (*user_entity.User, error) {
	workspaceMember, err := u.deps.WorkspaceMemberRepository.GetMemberById(ctx, memberId)
	if err != nil {
		return nil, err
	}
	user, err := u.deps.GetUserByIdUseCase.Execute(ctx, workspaceMember.UserId)
	if err != nil {
		return nil, err
	}
	return user, nil
}
