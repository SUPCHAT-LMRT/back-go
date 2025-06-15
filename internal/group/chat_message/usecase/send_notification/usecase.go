package send_notification

import (
	"context"
	entity2 "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/notification/entity"
	"github.com/supchat-lmrt/back-go/internal/notification/usecase/create_notification"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type SendMessageNotificationUseCaseDeps struct {
	uberdig.In
	Channels                  []Channel `group:"send_groupmessage_notification_channel"`
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
		Type:   entity.NotificationTypeGroupeMessage,
		IsRead: false,
		GroupMessageData: &entity.GroupMessageNotificationData{
			SenderId:  req.SenderId,
			GroupId:   req.GroupId,
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
	GroupId    group_entity.GroupId
	SenderId   user_entity.UserId
	MessageId  entity2.GroupChatMessageId
	ReceiverId user_entity.UserId
}
