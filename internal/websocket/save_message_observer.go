package websocket

import (
	"context"
	group_chat_message_entity "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	save_group_chat_message "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/save_message"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	save_channel_chat_message "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/save_message"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	uberdig "go.uber.org/dig"
	"time"
)

type SaveMessageObserverDeps struct {
	uberdig.In
	SaveChannelMessageUseCase   *save_channel_chat_message.SaveChannelMessageUseCase
	SaveGroupChatMessageUseCase *save_group_chat_message.SaveGroupChatMessageUseCase
	Logger                      logger.Logger
}

type SaveMessageObserver struct {
	deps SaveMessageObserverDeps
}

func NewSaveMessageObserver(deps SaveMessageObserverDeps) SendMessageObserver {
	return &SaveMessageObserver{deps: deps}
}

func (s SaveMessageObserver) OnSendMessage(message Message) {
	switch message.Target.Kind {
	case ChannelRoomKind:
		s.handleChannelMessage(message)
	case DirectRoomKind:
	// Todo
	//s.handleDirectMessage(message)
	case GroupRoomKind:
		s.handleGroupMessage(message)
	default:
		s.deps.Logger.Warn().Str("roomKind", string(message.Target.Kind)).Msg("Unknown room kind")
	}
}

func (s SaveMessageObserver) handleChannelMessage(message Message) {
	err := s.deps.SaveChannelMessageUseCase.Execute(context.Background(), &entity.ChannelMessage{
		ChannelId: channel_entity.ChannelId(message.Target.Id),
		Content:   message.Message,
		AuthorId:  message.Sender.UserId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save message")
	}
}

func (s SaveMessageObserver) handleGroupMessage(message Message) {
	err := s.deps.SaveGroupChatMessageUseCase.Execute(context.Background(), &group_chat_message_entity.GroupChatMessage{
		GroupId:   group_entity.GroupId(message.Target.Id),
		Content:   message.Message,
		AuthorId:  message.Sender.UserId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save message")
	}
}
