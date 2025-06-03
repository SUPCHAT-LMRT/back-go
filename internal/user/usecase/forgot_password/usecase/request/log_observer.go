package request

import (
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/entity"
)

type LogRequestForgotPasswordObserver struct {
	logger logger.Logger
}

func NewLogRequestForgotPasswordObserver(logger logger.Logger) ForgotPasswordRequestObserver {
	return &LogRequestForgotPasswordObserver{logger: logger}
}

func (o *LogRequestForgotPasswordObserver) NotifyRequestResetPasswordCreated(
	request entity.ForgotPasswordRequest,
) {
	o.logger.Info().Any("request", request).Msg("Forgot password request sent")
}
