package websocket

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/save_message"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"time"
)

type SaveMessageObserver struct {
	saveMessageUseCase *save_message.SaveMessageUseCase
	logger             logger.Logger
}

func NewSaveMessageObserver(saveMessageUseCase *save_message.SaveMessageUseCase, logger logger.Logger) SendMessageObserver {
	return &SaveMessageObserver{saveMessageUseCase: saveMessageUseCase, logger: logger}
}

func (s SaveMessageObserver) OnSendMessage(message Message) {
	// Todo, this is specific for channel message, found a way to use it for private/group messages
	err := s.saveMessageUseCase.Execute(context.Background(), &entity.ChannelMessage{
		ChannelId: channel_entity.ChannelId(message.Target.Id),
		Content:   message.Message,
		AuthorId:  message.Sender.UserId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.logger.Error().Err(err).Msg("Error on save message")
	}
}
