package request

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/service"
	uberdig "go.uber.org/dig"
)

type RequestForgotPasswordDeps struct {
	uberdig.In
	Service   service.ForgotPasswordService
	Observers []ForgotPasswordRequestObserver `group:"forgot_password_request_observers"`
}

type RequestForgotPasswordUseCase struct {
	deps RequestForgotPasswordDeps
}

func NewRequestForgotPasswordUseCase(deps RequestForgotPasswordDeps) *RequestForgotPasswordUseCase {
	return &RequestForgotPasswordUseCase{deps: deps}
}

func (u *RequestForgotPasswordUseCase) Execute(ctx context.Context, userId user_entity.UserId) (*entity.ForgotPasswordRequest, error) {
	forgotPasswordRequest, err := u.deps.Service.CreateForgotPasswordRequest(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, observer := range u.deps.Observers {
		go observer.NotifyRequestResetPasswordCreated(*forgotPasswordRequest)
	}

	return forgotPasswordRequest, err
}
