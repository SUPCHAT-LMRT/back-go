package reoder_channels

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"
	uberdig "go.uber.org/dig"
)

type ReorderChannelsInput struct {
	ChannelId entity.ChannelId
	NewIndex  int
}

type ReorderChannelsUseCaseDeps struct {
	uberdig.In
	Observers         []ReorderIndexChannelsObserver `group:"reorder_channels_observers"`
	Repo              repository.ChannelRepository
	GetChannelUseCase *get_channel.GetChannelUseCase
}

type ReorderChannelsUseCase struct {
	deps ReorderChannelsUseCaseDeps
}

func NewReorderChannelsUseCase(deps ReorderChannelsUseCaseDeps) *ReorderChannelsUseCase {
	return &ReorderChannelsUseCase{deps: deps}
}

func (uc *ReorderChannelsUseCase) ExecuteBulk(ctx context.Context, inputs []ReorderChannelsInput) error {

	if len(inputs) == 0 {
		return errors.New("no channels to reorder")
	}
	channel, err := uc.deps.GetChannelUseCase.Execute(ctx, inputs[0].ChannelId)
	if err != nil {
		return err
	}
	workspaceId := channel.WorkspaceId

	for _, input := range inputs {
		err := uc.deps.Repo.UpdateIndex(ctx, input.ChannelId, input.NewIndex)
		if err != nil {
			return err
		}
	}

	messages := make([]ChannelReorderMessage, len(inputs))

	for i, input := range inputs {
		messages[i] = ChannelReorderMessage{
			ChannelId: input.ChannelId,
			NewOrder:  input.NewIndex,
		}
	}

	for _, observer := range uc.deps.Observers {
		observer.NotifyChannelReordered(messages, workspaceId)
	}

	return nil
}
