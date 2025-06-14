package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	edit_message2 "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/edit_message"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/get_message"
	uberdig "go.uber.org/dig"
)

type SaveEditChannelMessageObserverDeps struct {
	uberdig.In
	EditChannelChatMessageUseCase *edit_message2.EditChannelChatMessageUseCase
	GetMessageUseCase             *get_message.GetMessageUseCase
	Logger                        logger.Logger
}

type SaveEditChannelMessageObserver struct {
	deps SaveEditChannelMessageObserverDeps
}

func NewSaveEditChannelMessageObserver(deps SaveEditChannelMessageObserverDeps) EditChannelMessageObserver {
	return &SaveEditChannelMessageObserver{deps: deps}
}

func (s SaveEditChannelMessageObserver) OnEditMessage(
	message *inbound.InboundChannelMessageContentEdit,
) {
	s.handleChannelMessage(message)
}

func (s SaveEditChannelMessageObserver) handleChannelMessage(
	message *inbound.InboundChannelMessageContentEdit,
) {
	// Fetch the original message to ensure it exists before editing
	originalMessage, err := s.deps.GetMessageUseCase.Execute(
		context.Background(),
		message.MessageId,
	)
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error fetching original group message")
		return
	}

	originalMessage.Content = message.NewContent

	err = s.deps.EditChannelChatMessageUseCase.Execute(
		context.Background(),
		originalMessage,
	)
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save group message")
	}
}
