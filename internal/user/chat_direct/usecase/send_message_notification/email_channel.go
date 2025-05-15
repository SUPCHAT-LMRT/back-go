package send_message_notification

import (
	"context"
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/mail"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	uberdig "go.uber.org/dig"
)

type EmailChannelDeps struct {
	uberdig.In
	GetUserByIdUseCase *get_by_id.GetUserByIdUseCase
	Mailer             *mail.Mailer
}

type EmailChannel struct {
	deps EmailChannelDeps
}

func NewEmailChannel(deps EmailChannelDeps) Channel {
	return &EmailChannel{deps: deps}
}

func (c *EmailChannel) SendNotification(ctx context.Context, req SendMessageNotificationRequest) error {
	receiver, err := c.deps.GetUserByIdUseCase.Execute(ctx, req.ReceiverId)
	if err != nil {
		return err
	}
	// get receiver notification preferences
	if !receiver.NotificationPreferences.Email {
		return nil
	}

	message := mail.NewMessage("Nouveau message privé", fmt.Sprintf("Vous avez reçu un message de %s avec le contenu %s", req.SenderName, req.Content))
	message.AddTo(receiver.Email)

	err = c.deps.Mailer.Send(message)
	if err != nil {
		return err
	}

	return nil
}
