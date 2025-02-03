package request

import "github.com/supchat-lmrt/back-go/internal/user/usecase/reset_password/entity"

type ResetPasswordRequestObserver interface {
	NotifyRequestResetPasswordCreated(request entity.ResetPasswordRequest)
}
