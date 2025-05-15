package send_message_notification

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	uberdig "go.uber.org/dig"
)

type SendMessageNotificationUseCaseDeps struct {
	uberdig.In
	Channels []Channel `group:"send_message_notification_channel"`
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

	return nil
}

type SendMessageNotificationRequest struct {
	Content    string
	SenderName string
	ReceiverId user_entity.UserId
}
