package request

import (
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/entity"
)

type LogRequestResetPasswordObserver struct {
	logger logger.Logger
}

func NewLogRequestResetPasswordObserver(logg logger.Logger) ResetPasswordRequestObserver {
	return &LogRequestResetPasswordObserver{logger: logg}
}

func (o *LogRequestResetPasswordObserver) NotifyRequestResetPasswordCreated(
	request entity.ResetPasswordRequest,
) {
	o.logger.Info().Any("request", request).Msg("Reset password request sent")
}
