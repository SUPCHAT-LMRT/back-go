package websocket

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/logger"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	save_channel_chat_message "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/save_message"
	uberdig "go.uber.org/dig"
)

type SaveChannelMessageObserverDeps struct {
	uberdig.In
	SaveChannelMessageUseCase *save_channel_chat_message.SaveChannelMessageUseCase
	Logger                    logger.Logger
}

type SaveChannelMessageObserver struct {
	deps SaveChannelMessageObserverDeps
}

func NewSaveChannelMessageObserver(deps SaveChannelMessageObserverDeps) SendChannelMessageObserver {
	return &SaveChannelMessageObserver{deps: deps}
}

func (s SaveChannelMessageObserver) OnSendMessage(
	message *inbound.InboundSendMessageToChannel,
	messageId entity.ChannelMessageId,
	userId user_entity.UserId,
) {
	s.handleChannelMessage(message, messageId, userId)
}

func (s SaveChannelMessageObserver) handleChannelMessage(
	message *inbound.InboundSendMessageToChannel,
	messageId entity.ChannelMessageId,
	userId user_entity.UserId,
) {
	err := s.deps.SaveChannelMessageUseCase.Execute(context.Background(), &entity.ChannelMessage{
		Id:        messageId,
		ChannelId: message.ChannelId,
		Content:   message.Content,
		AuthorId:  userId,
		CreatedAt: message.TransportMessageCreatedAt,
		UpdatedAt: message.TransportMessageCreatedAt,
	})
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save message")
	}
}
