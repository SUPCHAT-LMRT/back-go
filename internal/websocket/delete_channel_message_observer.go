package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/delete_message"
	uberdig "go.uber.org/dig"
)

type SaveDeleteChannelMessageObserverDeps struct {
	uberdig.In
	DeleteChannelChatMessageUseCase *delete_message.DeleteChannelChatMessageUseCase
	Logger                          logger.Logger
}

type SaveDeleteChannelMessageObserver struct {
	deps SaveDeleteChannelMessageObserverDeps
}

func NewSaveDeleteChannelMessageObserver(deps SaveDeleteChannelMessageObserverDeps) DeleteChannelMessageObserver {
	return &SaveDeleteChannelMessageObserver{deps: deps}
}

func (s SaveDeleteChannelMessageObserver) OnDeleteMessage(
	message *inbound.InboundChannelMessageDelete,
) {
	s.handleChannelMessage(message)
}

func (s SaveDeleteChannelMessageObserver) handleChannelMessage(
	message *inbound.InboundChannelMessageDelete,
) {
	err := s.deps.DeleteChannelChatMessageUseCase.Execute(
		context.Background(),
		message.MessageId,
	)
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save group message")
	}
}
