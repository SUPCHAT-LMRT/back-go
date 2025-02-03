package request

import (
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/entity"
)

type LogRequestValidationObserver struct {
	logger logger.Logger
}

func NewLogRequestValidationObserver(logger logger.Logger) ValidationRequestObserver {
	return &LogRequestValidationObserver{logger: logger}
}

func (o *LogRequestValidationObserver) NotifyRequestForgotPasswordCreated(request entity.ValidationRequest) {
	o.logger.Info().Any("request", request).Msg("Validation request sent")
}
