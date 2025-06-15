package send_notification

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	"github.com/supchat-lmrt/back-go/internal/notification/usecase/create_notification"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	channel_message_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	uberdig "go.uber.org/dig"
)

type SendMessageNotificationUseCaseDeps struct {
	uberdig.In
	Channels                  []Channel `group:"send_channelmessage_notification_channel"`
	CreateNotificationUseCase *create_notification.CreateNotificationUseCase
}

type SendMessageNotificationUseCase struct {
	deps SendMessageNotificationUseCaseDeps
}

func NewSendMessageNotificationUseCase(deps SendMessageNotificationUseCaseDeps) *SendMessageNotificationUseCase {
	return &SendMessageNotificationUseCase{deps: deps}
}

func (uc *SendMessageNotificationUseCase) Execute(ctx context.Context, req SendMessageNotificationRequest) error {
	for _, channel := range uc.deps.Channels {
		if err := channel.SendNotification(ctx, req); err != nil {
			return err
		}
	}

	notification := &entity.Notification{
		UserId: req.ReceiverId,
		Type:   entity.NotificationTypeChannelMessage,
		IsRead: false,
		ChannelMessageData: &entity.ChannelMessageNotificationData{
			SenderId:    req.SenderId,
			ChannelId:   req.ChannelId,
			WorkspaceId: req.WorkspaceId,
			MessageId:   req.MessageId,
		},
	}

	err := uc.deps.CreateNotificationUseCase.Execute(ctx, notification)
	if err != nil {
		return err
	}
	return nil
}

type SendMessageNotificationRequest struct {
	Content     string
	SenderName  string
	SenderId    user_entity.UserId
	MessageId   channel_message_entity.ChannelMessageId
	ReceiverId  user_entity.UserId
	ChannelId   channel_entity.ChannelId
	WorkspaceId workspace_entity.WorkspaceId
}
