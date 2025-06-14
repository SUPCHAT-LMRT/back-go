package group

import (
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"time"
)

type SearchGroup struct {
	Id        group_entity.GroupId
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}
