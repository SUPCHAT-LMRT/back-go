package request

import "github.com/supchat-lmrt/back-go/internal/user/usecase/validation/entity"

type ValidationRequestObserver interface {
	NotifyRequestForgotPasswordCreated(request entity.ValidationRequest)
}
