package create_channel

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/search/channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	uberdig "go.uber.org/dig"
)

type CreateChannelUseCaseDeps struct {
	uberdig.In
	Repository               repository.ChannelRepository
	SearchChannelSyncManager channel.SearchChannelSyncManager
	Observers                []CreateChannelObserver `group:"create_channel_observers"`
}

type CreateChannelUseCase struct {
	deps CreateChannelUseCaseDeps
}

func NewCreateChannelUseCase(deps CreateChannelUseCaseDeps) *CreateChannelUseCase {
	return &CreateChannelUseCase{deps: deps}
}

func (u *CreateChannelUseCase) Execute(ctx context.Context, chann *entity.Channel) error {
	err := u.deps.Repository.Create(ctx, chann)
	if err != nil {
		return err
	}

	err = u.deps.SearchChannelSyncManager.AddChannel(ctx, &channel.SearchChannel{
		Id:          chann.Id.String(),
		Name:        chann.Name,
		Topic:       chann.Topic,
		Kind:        channel.SearchChannelKindTextMessage,
		WorkspaceId: chann.WorkspaceId,
		CreatedAt:   chann.CreatedAt,
		UpdatedAt:   chann.CreatedAt,
	})
	if err != nil {
		return err
	}

	for _, observer := range u.deps.Observers {
		observer.ChannelCreated(chann)
	}

	return err
}
