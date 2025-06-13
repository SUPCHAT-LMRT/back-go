package event

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
)

const (
	GroupCreatedEventType EventType = "group_created"
)

type GroupCreatedEvent struct {
	Message *group_entity.Group
}

func (e GroupCreatedEvent) Type() EventType {
	return GroupCreatedEventType
}
