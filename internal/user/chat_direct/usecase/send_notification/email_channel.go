package send_notification

import (
	"context"
	"fmt"
	"log"

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

func (c *EmailChannel) SendNotification(
	ctx context.Context,
	req SendMessageNotificationRequest,
) error {
	receiver, err := c.deps.GetUserByIdUseCase.Execute(ctx, req.ReceiverId)
	if err != nil {
		log.Printf("Erreur lors de la récupération de l'utilisateur par son ID: %v", err)
		return err
	}
	// get receiver notification preferences
	// if !receiver.NotificationPreferences.Email {
	//	log.Printf("L'utilisateur %s a désactivé les notifications par email", receiver.Email)
	//	return nil
	// }

	message := mail.NewTextPlainMessage(
		"Nouveau message privé",
		fmt.Sprintf(
			"Vous avez reçu un message de %s avec le contenu %s",
			req.SenderName,
			req.Content,
		),
	)
	message.AddTo(receiver.Email)
	message.SetFrom(c.deps.Mailer.From)

	err = c.deps.Mailer.Send(message)
	if err != nil {
		log.Printf("Erreur lors de l'envoi du mail à %s: %v", receiver.Email, err)
		return err
	}
	log.Printf("L'email de notification à été envoyé à: %s", receiver.Email)
	return nil
}
