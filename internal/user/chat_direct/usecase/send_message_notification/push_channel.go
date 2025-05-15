package send_message_notification

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	uberdig "go.uber.org/dig"
)

type PushChannelDeps struct {
	uberdig.In
	GetUserByIdUseCase *get_by_id.GetUserByIdUseCase
}

type PushChannel struct {
	deps EmailChannelDeps
}

func NewPushChannel(deps EmailChannelDeps) Channel {
	return &PushChannel{deps: deps}
}

func (c *PushChannel) SendNotification(ctx context.Context, req SendMessageNotificationRequest) error {
	// get receiver notification preferences
	if !receiver.NotificationPreferences.Push {
		return nil
	}
	return nil
}
