package save_message

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/usecase/send_message_notification"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	uberdig "go.uber.org/dig"
)

type SendNotificationObserverDeps struct {
	uberdig.In
	GetUserByIdUseCase             *get_by_id.GetUserByIdUseCase
	SendMessageNotificationUseCase *send_message_notification.SendMessageNotificationUseCase
	Logger                         logger.Logger
}

type SendNotificationObserver struct {
	deps SendNotificationObserverDeps
}

func NewSendNotificationObserver(deps SendNotificationObserverDeps) MessageSavedObserver {
	return &SendNotificationObserver{deps: deps}
}

func (o SendNotificationObserver) NotifyMessageSaved(msg *entity.ChatDirect) {
	sender, err := o.deps.GetUserByIdUseCase.Execute(context.Background(), msg.SenderId)
	if err != nil {
		o.deps.Logger.Error().
			Str("chat_direct_id", msg.Id.String()).
			Err(err).Msg("failed to get sender user")
		return
	}

	err = o.deps.SendMessageNotificationUseCase.Execute(
		context.Background(),
		send_message_notification.SendMessageNotificationRequest{
			Content:    msg.Content,
			SenderName: sender.FullName(),
			ReceiverId: msg.GetReceiverId(),
		},
	)
	if err != nil {
		o.deps.Logger.Error().
			Str("chat_direct_id", msg.Id.String()).
			Str("sender_name", sender.FullName()).
			Str("receiver_id", msg.GetReceiverId().String()).
			Err(err).Msg("failed to send message notification")
		return
	}
}
