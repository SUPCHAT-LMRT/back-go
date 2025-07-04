package sendmail

import (
	"github.com/supchat-lmrt/back-go/internal/mail"
)

type SendMailUseCase struct {
	mailer *mail.Mailer
}

func NewSendMailUseCase(mailer *mail.Mailer) *SendMailUseCase {
	return &SendMailUseCase{mailer: mailer}
}

func (u *SendMailUseCase) Execute(message *mail.Message) error {
	if message.From == "" {
		message.From = u.mailer.From
	}
	return u.mailer.Send(message)
}

func (u *SendMailUseCase) Outros() []string {
	return []string{
		"Besoin d'aide, ou avez-vous des questions? Vous pouvez répondre à cet e-mail pour nous contacter.",
		"A bientôt sur Supchat-LMRT!",
	}
}
