package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/edit_message"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/get_message"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	uberdig "go.uber.org/dig"
)

type SaveEditGroupMessageObserverDeps struct {
	uberdig.In
	EditGroupChatMessageUseCase *edit_message.EditGroupChatMessageUseCase
	GetMessageUseCase           *get_message.GetMessageUseCase
	Logger                      logger.Logger
}

type SaveEditGroupMessageObserver struct {
	deps SaveEditGroupMessageObserverDeps
}

func NewSaveEditGroupMessageObserver(deps SaveEditGroupMessageObserverDeps) EditGroupMessageObserver {
	return &SaveEditGroupMessageObserver{deps: deps}
}

func (s SaveEditGroupMessageObserver) OnEditMessage(
	message *inbound.InboundGroupMessageContentEdit,
) {
	s.handleGroupMessage(message)
}

func (s SaveEditGroupMessageObserver) handleGroupMessage(
	message *inbound.InboundGroupMessageContentEdit,
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

	err = s.deps.EditGroupChatMessageUseCase.Execute(
		context.Background(),
		originalMessage,
	)
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save group message")
	}
}
