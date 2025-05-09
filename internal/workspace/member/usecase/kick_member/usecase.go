package kick_member

import (
	"context"
	"errors"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity3 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
)

var ErrMemberNotFound = errors.New("member not found in workspace")

type KickMemberUseCase struct {
	WorkspaceMemberRepository repository.WorkspaceMemberRepository
}

func NewKickMemberUseCase(repo repository.WorkspaceMemberRepository) *KickMemberUseCase {
	return &KickMemberUseCase{WorkspaceMemberRepository: repo}
}

func (u *KickMemberUseCase) Execute(ctx context.Context, workspaceId entity2.WorkspaceId, memberId entity3.WorkspaceMemberId) error {
	exists, err := u.WorkspaceMemberRepository.IsMemberExists(ctx, workspaceId, memberId)
	if err != nil {
		return err
	}
	if !exists {
		return ErrMemberNotFound
	}

	return u.WorkspaceMemberRepository.RemoveMember(ctx, workspaceId, memberId)
}
