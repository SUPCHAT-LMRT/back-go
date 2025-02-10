package request

import (
	"context"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/service"
	uberdig "go.uber.org/dig"
)

type RequestResetPasswordDeps struct {
	uberdig.In
	Service   service.ResetPasswordService
	Observers []ResetPasswordRequestObserver `group:"reset_password_request_observers"`
}

type RequestResetPasswordUseCase struct {
	deps RequestResetPasswordDeps
}

func NewRequestResetPasswordUseCase(deps RequestResetPasswordDeps) *RequestResetPasswordUseCase {
	return &RequestResetPasswordUseCase{deps: deps}
}

func (u *RequestResetPasswordUseCase) Execute(ctx context.Context, userId user_entity.UserId) (*entity.ResetPasswordRequest, error) {
	resetPasswordRequest, err := u.deps.Service.CreateResetPasswordRequest(ctx, userId)
	if err != nil {
		return nil, err
	}

	for _, observer := range u.deps.Observers {
		observer.NotifyRequestResetPasswordCreated(*resetPasswordRequest)
	}

	return resetPasswordRequest, err
}
