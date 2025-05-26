package validate

import (
	"context"

	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/service"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/update_password"
)

type ValidateForgotPasswordUseCase struct {
	service               service.ForgotPasswordService
	updatePasswordUseCase *update_password.ChangePasswordUseCase
}

func NewValidateForgotPasswordUseCase(
	service service.ForgotPasswordService,
	updatePasswordUseCase *update_password.ChangePasswordUseCase,
) *ValidateForgotPasswordUseCase {
	return &ValidateForgotPasswordUseCase{
		service:               service,
		updatePasswordUseCase: updatePasswordUseCase,
	}
}

func (u *ValidateForgotPasswordUseCase) Execute(
	ctx context.Context,
	token uuid.UUID,
	password string,
) error {
	user, err := u.service.DeleteForgotPasswordRequest(ctx, token)
	if err != nil {
		return err
	}

	// TODO: call observers here (ex: SendEmailObserver)

	err = u.updatePasswordUseCase.Execute(ctx, user, password)
	if err != nil {
		return err
	}

	return nil
}
