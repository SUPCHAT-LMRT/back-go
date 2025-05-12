package generate

import (
	"github.com/matcornic/hermes/v2"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/mail"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/sendmail"
)

type SendMailGenerateInviteLinkObserver struct {
	logger          logger.Logger
	sendMailUseCase *sendmail.SendMailUseCase
}

func NewSendMailGenerateInviteLinkObserver(logger logger.Logger, sendMailUseCase *sendmail.SendMailUseCase) GenerateInviteLinkObserver {
	return &SendMailGenerateInviteLinkObserver{logger: logger, sendMailUseCase: sendMailUseCase}
}

func (o *SendMailGenerateInviteLinkObserver) NotifyInviteLinkGenerated(inviteLink *entity.InviteLink, link string) {
	outros := o.sendMailUseCase.Outros()
	outros = append(outros, "Si vous n'êtes pas à l'origine de cette demande, veuillez ignorer ce message.")

	email := hermes.Email{
		Body: hermes.Body{
			Greeting:  "Bonjour",
			Signature: "Cordialement",
			Name:      inviteLink.FirstName + " " + inviteLink.LastName,
			Intros: []string{
				"Vous avez été invité à rejoindre Supchat-LMRT.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Cliquez sur le bouton ci-dessous pour vous inscrire.",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "S'inscrire",
						Link:  link,
					},
				},
			},
			Outros: outros,
		},
	}

	h := hermes.Hermes{
		// Optional Theme
		Theme: new(hermes.Default),
		Product: hermes.Product{
			// Appears in header & footer of e-mails
			Name:        "Supchat-LMRT",
			Link:        "https://supchart-lmrt.fr",
			Copyright:   "© 2024 Supchat-LMRT",
			TroubleText: "Si vous rencontrez des problèmes en cliquant sur le bouton '{ACTION}', copiez et collez l'URL ci-dessous dans votre navigateur Web:",
		},
	}
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		o.logger.Error().Err(err).Msg("Error generating email body")
		return
	}
	msg := mail.NewHTMLMessage("Inscription à Supchat-LMRT", emailBody)
	msg.AddTo(inviteLink.Email)

	err = o.sendMailUseCase.Execute(msg)
	if err != nil {
		o.logger.Error().Err(err).Msg("Error sending email")
		return
	}
}
