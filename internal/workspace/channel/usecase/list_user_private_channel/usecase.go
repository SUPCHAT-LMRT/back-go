package list_user_private_channel

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	workspace_member_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
)

type ListPrivateChannelMembersUseCase struct {
	repo repository.ChannelRepository
}

func NewListPrivateChannelMembersUseCase(
	repo repository.ChannelRepository,
) *ListPrivateChannelMembersUseCase {
	return &ListPrivateChannelMembersUseCase{repo: repo}
}

func (uc *ListPrivateChannelMembersUseCase) Execute(
	ctx context.Context,
	channelId entity.ChannelId,
) ([]workspace_member_entity.WorkspaceMemberId, error) {
	return uc.repo.ListMembersOfPrivateChannel(ctx, channelId)
}
