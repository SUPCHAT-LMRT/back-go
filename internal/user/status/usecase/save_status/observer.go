package save_status

import user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"

type UserStatusSavedObserver interface {
	NotifyUserStatusSaved(userStatus *user_status_entity.UserStatus)
}
