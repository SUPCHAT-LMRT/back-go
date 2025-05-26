package add_member

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/repository"
)

type AddMemberUseCase struct {
	workspaceMemberRepository repository.WorkspaceMemberRepository
}

func NewAddMemberUseCase(
	workspaceMemberRepository repository.WorkspaceMemberRepository,
) *AddMemberUseCase {
	return &AddMemberUseCase{workspaceMemberRepository: workspaceMemberRepository}
}

func (u *AddMemberUseCase) Execute(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	member *entity2.WorkspaceMember,
) error {
	return u.workspaceMemberRepository.AddMember(ctx, workspaceId, member)
}
