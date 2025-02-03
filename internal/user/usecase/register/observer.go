package register

import "github.com/supchat-lmrt/back-go/internal/user/entity"

type RegisterUserObserver interface {
	NotifyUserRegistered(user entity.User)
}
