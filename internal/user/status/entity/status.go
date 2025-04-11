package entity

import user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"

type Status string

const (
	StatusUnknown      Status = "unknown"
	StatusOnline       Status = "online"
	StatusDoNotDisturb Status = "do-not-disturb"
	StatusInvisible    Status = "invisible"
	StatusAway         Status = "away"
	StatusOffline      Status = "offline"
)

type UserStatus struct {
	UserId user_entity.UserId
	Status Status
}

func (s Status) String() string {
	return string(s)
}

func ParseStatus(status string) Status {
	switch status {
	case StatusOnline.String():
		return StatusOnline
	case StatusDoNotDisturb.String():
		return StatusDoNotDisturb
	case StatusInvisible.String():
		return StatusInvisible
	case StatusAway.String():
		return StatusAway
	case StatusOffline.String():
		return StatusOffline
	default:
		return StatusUnknown
	}
}
