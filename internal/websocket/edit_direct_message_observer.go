package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	edit_message2 "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/edit_message"
	get_message2 "github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/get_message"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	uberdig "go.uber.org/dig"
)

type SaveEditDirectMessageObserverDeps struct {
	uberdig.In
	EditDirectChatMessageUseCase *edit_message2.EditDirectChatMessageUseCase
	GetDirectChatMessageUseCase  *get_message2.GetDirectChatMessageUseCase
	Logger                       logger.Logger
}

type SaveEditDirectMessageObserver struct {
	deps SaveEditDirectMessageObserverDeps
}

func NewSaveEditDirectMessageObserver(deps SaveEditDirectMessageObserverDeps) EditDirectMessageObserver {
	return &SaveEditDirectMessageObserver{deps: deps}
}

func (s SaveEditDirectMessageObserver) OnEditMessage(
	message *inbound.InboundDirectMessageContentEdit,
) {
	s.handleDirectMessage(message)
}

func (s SaveEditDirectMessageObserver) handleDirectMessage(
	message *inbound.InboundDirectMessageContentEdit,
) {
	// Fetch the original message to ensure it exists before editing
	originalMessage, err := s.deps.GetDirectChatMessageUseCase.Execute(
		context.Background(),
		message.MessageId,
	)
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error fetching original group message")
		return
	}

	originalMessage.Content = message.NewContent

	err = s.deps.EditDirectChatMessageUseCase.Execute(
		context.Background(),
		originalMessage,
	)
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save group message")
	}
}
