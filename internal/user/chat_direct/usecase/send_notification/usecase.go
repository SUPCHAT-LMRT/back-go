package send_notification

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	"github.com/supchat-lmrt/back-go/internal/notification/usecase/create_notification"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type SendMessageNotificationUseCaseDeps struct {
	uberdig.In
	Channels                  []Channel `group:"send_directmessage_notification_channel"`
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
		Type:   entity.NotificationTypeDirectMessage,
		IsRead: false,
		DirectMessageData: &entity.DirectMessageNotificationData{
			SenderId:  req.SenderId,
			MessageId: req.MessageId,
		},
	}

	err := uc.deps.CreateNotificationUseCase.Execute(ctx, notification)
	if err != nil {
		return err
	}
	return nil
}

type SendMessageNotificationRequest struct {
	Content    string
	SenderName string
	SenderId   user_entity.UserId
	MessageId  chat_direct_entity.ChatDirectId
	ReceiverId user_entity.UserId
}
