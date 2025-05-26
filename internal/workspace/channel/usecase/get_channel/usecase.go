package get_channel

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
)

type GetChannelUseCase struct {
	repository repository.ChannelRepository
}

func NewGetChannelUseCase(channelRepository repository.ChannelRepository) *GetChannelUseCase {
	return &GetChannelUseCase{repository: channelRepository}
}

func (u *GetChannelUseCase) Execute(
	ctx context.Context,
	channelId entity.ChannelId,
) (*entity.Channel, error) {
	return u.repository.GetById(ctx, channelId)
}
