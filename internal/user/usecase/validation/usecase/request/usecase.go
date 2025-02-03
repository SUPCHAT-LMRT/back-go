package request

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/service"
	uberdig "go.uber.org/dig"
)

type RequestAccountValidationDeps struct {
	uberdig.In
	Service   service.ValidationService
	Observers []ValidationRequestObserver `group:"validation_request_observers"`
}

type RequestAccountValidationUseCase struct {
	deps RequestAccountValidationDeps
}

func NewRequestAccountValidationUseCase(deps RequestAccountValidationDeps) *RequestAccountValidationUseCase {
	return &RequestAccountValidationUseCase{deps: deps}
}

func (u *RequestAccountValidationUseCase) Execute(ctx context.Context, userId user_entity.UserId) (*entity.ValidationRequest, error) {
	validationRequest, err := u.deps.Service.CreateAccountValidationRequest(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, observer := range u.deps.Observers {
		go observer.NotifyRequestForgotPasswordCreated(*validationRequest)
	}

	return validationRequest, err
}
