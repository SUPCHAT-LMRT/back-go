package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/save_message"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	uberdig "go.uber.org/dig"
	"time"
)

type SaveDirectMessageObserverDeps struct {
	uberdig.In
	SaveDirectMessageUseCase *save_message.SaveDirectMessageUseCase
	Logger                   logger.Logger
}

type SaveDirectMessageObserver struct {
	deps SaveDirectMessageObserverDeps
}

func NewSaveDirectMessageObserver(deps SaveDirectMessageObserverDeps) SendDirectMessageObserver {
	return &SaveDirectMessageObserver{deps: deps}
}

func (s SaveDirectMessageObserver) OnSendMessage(message *inbound.InboundSendDirectMessage, messageId chat_direct_entity.ChatDirectId, userId user_entity.UserId) {
	s.handleDirectMessage(message, messageId, userId)
}

func (s SaveDirectMessageObserver) handleDirectMessage(message *inbound.InboundSendDirectMessage, messageId chat_direct_entity.ChatDirectId, userId user_entity.UserId) {
	now := time.Now()
	err := s.deps.SaveDirectMessageUseCase.Execute(context.Background(), &chat_direct_entity.ChatDirect{
		Id:        messageId,
		SenderId:  userId,
		User1Id:   userId,
		User2Id:   message.OtherUserId,
		Content:   message.Content,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save message")
	}
}
