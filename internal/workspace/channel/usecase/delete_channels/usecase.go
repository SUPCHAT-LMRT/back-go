package delete_channels

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/search/channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	uberdig "go.uber.org/dig"
)

type DeleteChannelUseCaseDeps struct {
	uberdig.In
	Repository               repository.ChannelRepository
	SearchChannelSyncManager channel.SearchChannelSyncManager
	Observers                []DeleteSpecifyChannelsObserver `group:"delete_channels_observers"`
}

type DeleteChannelUseCase struct {
	deps DeleteChannelUseCaseDeps
}

func NewDeleteChannelUseCase(deps DeleteChannelUseCaseDeps) *DeleteChannelUseCase {
	return &DeleteChannelUseCase{deps: deps}
}

func (u *DeleteChannelUseCase) Execute(ctx context.Context, channelId entity.ChannelId) error {
	id, err2 := u.deps.Repository.GetById(ctx, channelId)
	if err2 != nil {
		return err2
	}

	err := u.deps.Repository.Delete(ctx, channelId)
	if err != nil {
		return err
	}

	err = u.deps.SearchChannelSyncManager.RemoveChannel(ctx, channelId)
	if err != nil {
		return err
	}

	for _, observer := range u.deps.Observers {
		observer.NotifyChannelsDeleted(channelId, id.WorkspaceId)
	}

	return nil
}
