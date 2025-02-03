package register

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/usecase/request"
)

type RequestValidationObserver struct {
	logger  logger.Logger
	useCase *request.RequestAccountValidationUseCase
}

func NewRequestValidationObserver(logger logger.Logger, requestAccountValidationUseCase *request.RequestAccountValidationUseCase) RegisterUserObserver {
	return &RequestValidationObserver{logger: logger, useCase: requestAccountValidationUseCase}
}

func (o *RequestValidationObserver) NotifyUserRegistered(user entity.User) {
	_, err := o.useCase.Execute(context.Background(), user.Id)
	if err != nil {
		o.logger.Error().Err(err).Send()
		return
	}
}
