package request

import (
	"os"
	"strings"

	"github.com/matcornic/hermes/v2"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/mail"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/sendmail"
)

type SendMailRequestResetPasswordObserver struct {
	logger          logger.Logger
	sendMailUseCase *sendmail.SendMailUseCase
}

func NewSendMailRequestResetPasswordObserver(
	logger logger.Logger,
	sendMailUseCase *sendmail.SendMailUseCase,
) ResetPasswordRequestObserver {
	return &SendMailRequestResetPasswordObserver{logger: logger, sendMailUseCase: sendMailUseCase}
}

func (o *SendMailRequestResetPasswordObserver) NotifyRequestResetPasswordCreated(
	request entity.ResetPasswordRequest,
) {
	user := request.User

	validateUrl := os.Getenv("FRONT_ACCOUNT_RESET_PASSWORD_URL")
	if validateUrl == "" {
		o.logger.Warn().Msg("FRONT_ACCOUNT_RESET_PASSWORD_URL is not set")
		return
	}

	validateUrl = strings.Replace(validateUrl, "{token}", request.Token.String(), 1)

	outros := o.sendMailUseCase.Outros()
	outros = append(
		outros,
		"Si vous n'êtes pas à l'origine de cette demande, veuillez ignorer ce message.",
	)

	email := hermes.Email{
		Body: hermes.Body{
			Greeting:  "Bonjour",
			Signature: "Cordialement",
			Name:      user.FullName(),
			Intros: []string{
				"Nous avons reçu une demande de réinitialisation de mot de passe pour votre compte Supchat-LMRT.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Pour réinitialiser votre mot de passe, veuillez cliquer sur le bouton ci-dessous:",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Réinitialiser mon mot de passe",
						Link:  validateUrl,
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
			Link:        "https://supchat-lmrt.fr",
			Copyright:   "© 2025 Supchat-LMRT",
			TroubleText: "Si vous rencontrez des problèmes en cliquant sur le bouton '{ACTION}', copiez et collez l'URL ci-dessous dans votre navigateur Web:",
		},
	}
	emailBody, err := h.GenerateHTML(email)
	if err != nil {
		o.logger.Error().Err(err).Msg("Error generating email body")
		return
	}
	msg := mail.NewHTMLMessage("Réinitialisation du mot de passe", emailBody)
	msg.AddTo(user.Email)

	err = o.sendMailUseCase.Execute(msg)
	if err != nil {
		o.logger.Error().Err(err).Msg("Error sending email")
		return
	}
}
