package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/delete_message"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	uberdig "go.uber.org/dig"
)

type SaveDeleteGroupMessageObserverDeps struct {
	uberdig.In
	DeleteGroupChatMessageUseCase *delete_message.DeleteGroupChatMessageUseCase
	Logger                        logger.Logger
}

type SaveDeleteGroupMessageObserver struct {
	deps SaveDeleteGroupMessageObserverDeps
}

func NewSaveDeleteGroupMessageObserver(deps SaveDeleteGroupMessageObserverDeps) DeleteGroupMessageObserver {
	return &SaveDeleteGroupMessageObserver{deps: deps}
}

func (s SaveDeleteGroupMessageObserver) OnDeleteMessage(
	message *inbound.InboundGroupMessageDelete,
) {
	s.handleGroupMessage(message)
}

func (s SaveDeleteGroupMessageObserver) handleGroupMessage(
	message *inbound.InboundGroupMessageDelete,
) {
	err := s.deps.DeleteGroupChatMessageUseCase.Execute(
		context.Background(),
		message.MessageId,
	)
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save group message")
	}
}
