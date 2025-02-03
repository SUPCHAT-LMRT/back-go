package request

import (
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/entity"
)

type LogRequestResetPasswordObserver struct {
	logger logger.Logger
}

func NewLogRequestResetPasswordObserver(logger logger.Logger) ResetPasswordRequestObserver {
	return &LogRequestResetPasswordObserver{logger: logger}
}

func (o *LogRequestResetPasswordObserver) NotifyRequestResetPasswordCreated(request entity.ResetPasswordRequest) {
	o.logger.Info().Any("request", request).Msg("Reset password request sent")
}
