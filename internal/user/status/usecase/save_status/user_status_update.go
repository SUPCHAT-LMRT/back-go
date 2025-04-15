package save_status

import (
	"github.com/supchat-lmrt/back-go/internal/event"
	"github.com/supchat-lmrt/back-go/internal/logger"
	user_status_entity "github.com/supchat-lmrt/back-go/internal/user/status/entity"
	uberdig "go.uber.org/dig"
)

type UserStatusUpdateObserverDeps struct {
	uberdig.In
	EventBus *event.EventBus
	Logger   logger.Logger
}

type UserStatusUpdateObserver struct {
	deps UserStatusUpdateObserverDeps
}

func NewUserStatusUpdateObserver(deps UserStatusUpdateObserverDeps) UserStatusSavedObserver {
	return &UserStatusUpdateObserver{deps: deps}
}

func (o UserStatusUpdateObserver) NotifyUserStatusSaved(userStatus *user_status_entity.UserStatus) {
	// Publish an event after saving the message
	o.deps.EventBus.Publish(&event.UserStatusSavedEvent{
		UserStatus: userStatus,
	})
}
