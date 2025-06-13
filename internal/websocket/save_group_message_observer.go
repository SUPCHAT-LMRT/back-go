package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/save_message"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/logger"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	uberdig "go.uber.org/dig"
)

type SaveGroupMessageObserverDeps struct {
	uberdig.In
	SaveGroupChatMessageUseCase *save_message.SaveGroupChatMessageUseCase
	Logger                      logger.Logger
}

type SaveGroupMessageObserver struct {
	deps SaveGroupMessageObserverDeps
}

func NewSaveGroupMessageObserver(deps SaveGroupMessageObserverDeps) SendGroupMessageObserver {
	return &SaveGroupMessageObserver{deps: deps}
}

func (s SaveGroupMessageObserver) OnSendMessage(
	message *inbound.InboundSendGroupMessage,
	messageId entity.GroupChatMessageId,
	userId user_entity.UserId,
) {
	s.handleGroupMessage(message, messageId, userId)
}

func (s SaveGroupMessageObserver) handleGroupMessage(
	message *inbound.InboundSendGroupMessage,
	messageId entity.GroupChatMessageId,
	userId user_entity.UserId,
) {
	err := s.deps.SaveGroupChatMessageUseCase.Execute(
		context.Background(),
		&entity.GroupChatMessage{
			Id:        messageId,
			GroupId:   group_entity.GroupId(message.GroupId),
			Content:   message.Content,
			AuthorId:  userId,
			CreatedAt: message.TransportMessageCreatedAt,
			UpdatedAt: message.TransportMessageCreatedAt,
		},
	)
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save group message")
	}
}
