package create_channel

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	uberdig "go.uber.org/dig"
)

type CreateChannelUseCaseDeps struct {
	uberdig.In
	Repository repository.ChannelRepository
	Observers  []CreateChannelObserver `group:"create_channel_observers"`
}

type CreateChannelUseCase struct {
	deps CreateChannelUseCaseDeps
}

func NewCreateChannelUseCase(deps CreateChannelUseCaseDeps) *CreateChannelUseCase {
	return &CreateChannelUseCase{deps: deps}
}

func (u *CreateChannelUseCase) Execute(ctx context.Context, channel *entity.Channel) error {
	err := u.deps.Repository.Create(ctx, channel)
	if err != nil {
		return err
	}

	for _, observer := range u.deps.Observers {
		observer.ChannelCreated(channel)
	}

	return err
}
