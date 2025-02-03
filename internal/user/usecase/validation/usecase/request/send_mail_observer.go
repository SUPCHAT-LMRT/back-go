package request

import (
	"github.com/matcornic/hermes/v2"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/mail"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/sendmail"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/entity"
	"os"
	"strings"
)

type SendMailRequestValidationObserver struct {
	logger          logger.Logger
	sendMailUseCase *sendmail.SendMailUseCase
}

func NewSendMailRequestValidationObserver(logger logger.Logger, sendMailUseCase *sendmail.SendMailUseCase) ValidationRequestObserver {
	return &SendMailRequestValidationObserver{logger: logger, sendMailUseCase: sendMailUseCase}
}

func (o *SendMailRequestValidationObserver) NotifyRequestForgotPasswordCreated(request entity.ValidationRequest) {
	user := request.User

	validateUrl := os.Getenv("FRONT_ACCOUNT_VALIDATE_URL")
	if validateUrl == "" {
		o.logger.Warn().Msg("FRONT_ACCOUNT_VALIDATE_URL is not set")
		return
	}

	validateUrl = strings.Replace(validateUrl, "{token}", request.Token.String(), 1)

	email := hermes.Email{
		Body: hermes.Body{
			Greeting:  "Bonjour",
			Signature: "Cordialement",
			Name:      user.FirstName + " " + user.LastName,
			Intros: []string{
				"Bienvenue sur Supchat-LMRT, nous vous remercions pour votre inscription.",
			},
			Actions: []hermes.Action{
				{
					Instructions: "Pour commencer à utiliser Supchat-LMRT, veuillez confirmer votre compte en cliquant sur le bouton ci-dessous:",
					Button: hermes.Button{
						Color: "#22BC66", // Optional action button color
						Text:  "Confirmer mon compte",
						Link:  validateUrl,
					},
				},
			},
			Outros: o.sendMailUseCase.Outros(),
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
	msg := mail.NewHTMLMessage("Bienvenue sur Supchat-LMRT", emailBody)
	msg.AddTo(user.Email)

	err = o.sendMailUseCase.Execute(msg)
	if err != nil {
		o.logger.Error().Err(err).Msg("Error sending email")
		return
	}

}
