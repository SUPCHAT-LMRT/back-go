package websocket

import (
	"context"
	save_group_chat_message "github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/save_message"
	"github.com/supchat-lmrt/back-go/internal/logger"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	save_channel_chat_message "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/save_message"
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

func (s SaveMessageObserver) OnSendMessage(message messages.Message, messageId entity.ChannelMessageId, userId user_entity.UserId) {
	switch msg := message.(type) {
	case *inbound.InboundSendMessageToChannel:
		s.handleChannelMessage(msg, messageId, userId)
	}
	//switch message.Target.Kind {
	//case ChannelRoomKind:
	//	s.handleChannelMessage(message)
	//case DirectRoomKind:
	//// Todo
	////s.handleDirectMessage(message)
	//case GroupRoomKind:
	//	s.handleGroupMessage(message)
	//default:
	//	s.deps.Logger.Warn().Str("roomKind", string(message.Target.Kind)).Msg("Unknown room kind")
	//}
}

func (s SaveMessageObserver) handleChannelMessage(message *inbound.InboundSendMessageToChannel, messageId entity.ChannelMessageId, userId user_entity.UserId) {
	err := s.deps.SaveChannelMessageUseCase.Execute(context.Background(), &entity.ChannelMessage{
		Id:        messageId,
		ChannelId: message.ChannelId,
		Content:   message.Content,
		AuthorId:  userId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		s.deps.Logger.Error().Err(err).Msg("Error on save message")
	}
}

func (s SaveMessageObserver) handleGroupMessage(message messages.Message) {
	//err := s.deps.SaveGroupChatMessageUseCase.Execute(context.Background(), &group_chat_message_entity.GroupChatMessage{
	//	GroupId:   group_entity.GroupId(message.Target.ChannelId),
	//	Content:   message.Message,
	//	AuthorId:  message.Sender.UserId,
	//	CreatedAt: time.Now(),
	//})
	//if err != nil {
	//	s.deps.Logger.Error().Err(err).Msg("Error on save message")
	//}
}
