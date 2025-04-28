package create_channel

import (
	"context"
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/search/channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/repository"
	uberdig "go.uber.org/dig"
	"time"
)

type CreateChannelUseCaseDeps struct {
	uberdig.In
	Repository               repository.ChannelRepository
	SearchChannelSyncManager channel.SearchChannelSyncManager
	Observers                []CreateSpecifyChannelObserver `group:"create_channel_observers"`
}

type CreateChannelUseCase struct {
	deps CreateChannelUseCaseDeps
}

func NewCreateChannelUseCase(deps CreateChannelUseCaseDeps) *CreateChannelUseCase {
	return &CreateChannelUseCase{deps: deps}
}

func (u *CreateChannelUseCase) Execute(ctx context.Context, chann *entity.Channel) error {
	channelCount, err := u.deps.Repository.CountByWorkspaceId(ctx, chann.WorkspaceId)
	if err != nil {
		return err
	}
	chann.CreatedAt = time.Now()
	chann.UpdatedAt = chann.CreatedAt
	chann.Index = int(channelCount)

	err = u.deps.Repository.Create(ctx, chann)
	if err != nil {
		return err
	}

	err = u.deps.SearchChannelSyncManager.AddChannel(ctx, &channel.SearchChannel{
		Id:          chann.Id,
		Name:        chann.Name,
		Topic:       chann.Topic,
		Kind:        mapChannelKindToSearchResultChannelKind(chann.Kind),
		WorkspaceId: chann.WorkspaceId,
		CreatedAt:   chann.CreatedAt,
		UpdatedAt:   chann.UpdatedAt,
	})
	if err != nil {
		return err
	}

	fmt.Println(u.deps.Observers)

	for _, observer := range u.deps.Observers {
		observer.NotifyChannelCreated(chann)
	}

	return err
}

func mapChannelKindToSearchResultChannelKind(kind entity.ChannelKind) channel.SearchChannelKind {
	switch kind {
	case entity.ChannelKindText:
		return channel.SearchChannelKindText
	case entity.ChannelKindVoice:
		return channel.SearchChannelKindVoice
	default:
		return channel.SearchChannelKindUnknown
	}
}
