package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/delete_message"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	uberdig "go.uber.org/dig"
)

type SaveDeleteDirectMessageObserverDeps struct {
	uberdig.In
	DeleteDirectChatMessageUseCase *delete_message.DeleteDirectChatMessageUseCase
	Logger                         logger.Logger
}

type SaveDeleteDirectMessageObserver struct {
	deps SaveDeleteDirectMessageObserverDeps
}

func NewSaveDeleteDirectMessageObserver(deps SaveDeleteDirectMessageObserverDeps) DeleteDirectMessageObserver {
	return &SaveDeleteDirectMessageObserver{deps: deps}
}

func (s SaveDeleteDirectMessageObserver) OnDeleteMessage(
	message *inbound.InboundDirectMessageDelete,
) {
	s.handleDirectMessage(message)
}

func (s SaveDeleteDirectMessageObserver) handleDirectMessage(
	message *inbound.InboundDirectMessageDelete,
) {
	err := s.deps.DeleteDirectChatMessageUseCase.Execute(
		context.Background(),
		message.MessageId,
	)
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save group message")
	}
}
