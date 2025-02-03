package validate

import (
	"context"
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/service"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/update_password"
)

type ValidateResetPasswordUseCase struct {
	service               service.ResetPasswordService
	updatePasswordUseCase *update_password.ChangePasswordUseCase
}

func NewValidateResetPasswordUseCase(service service.ResetPasswordService, updatePasswordUseCase *update_password.ChangePasswordUseCase) *ValidateResetPasswordUseCase {
	return &ValidateResetPasswordUseCase{service: service, updatePasswordUseCase: updatePasswordUseCase}
}

func (u *ValidateResetPasswordUseCase) Execute(ctx context.Context, token uuid.UUID, password string) error {
	user, err := u.service.DeleteResetPasswordRequest(ctx, token)
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
