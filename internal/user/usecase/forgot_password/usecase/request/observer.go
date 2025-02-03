package request

import "github.com/supchat-lmrt/back-go/internal/user/usecase/forgot_password/entity"

type ForgotPasswordRequestObserver interface {
	NotifyRequestResetPasswordCreated(request entity.ForgotPasswordRequest)
}
