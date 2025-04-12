package event

import user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"

const (
	UserStatusSavedEventType EventType = "user_status_saved"
)

type UserStatusSavedEvent struct {
	UserStatus *user_status_entity.UserStatus
}

func (e UserStatusSavedEvent) Type() EventType {
	return UserStatusSavedEventType
}
