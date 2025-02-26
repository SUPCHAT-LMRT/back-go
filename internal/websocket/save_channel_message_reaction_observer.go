package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/add_reaction"
	uberdig "go.uber.org/dig"
)

type SaveChannelMessageReactionObserverDeps struct {
	uberdig.In
	AddReactionUseCase *add_reaction.AddReactionUseCase
	Logger             logger.Logger
}

type SaveChannelMessageReactionObserver struct {
	deps SaveChannelMessageReactionObserverDeps
}

func NewSaveChannelMessageReactionObserver(deps SaveChannelMessageReactionObserverDeps) CreateChannelMessageReactionObserver {
	return &SaveChannelMessageReactionObserver{deps: deps}
}

func (s SaveChannelMessageReactionObserver) OnCreateChannelMessageReaction(messageId entity.ChannelMessageId, reactionId entity.ChannelMessageReactionId, userId user_entity.UserId, reaction string) {
	err := s.deps.AddReactionUseCase.Execute(context.Background(), entity.ChannelMessageReaction{
		Id:        reactionId,
		MessageId: messageId,
		Reaction:  reaction,
		UserId:    userId,
	})
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("failed to save reaction")
	}
}
