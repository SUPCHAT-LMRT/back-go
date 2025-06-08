package entity

import (
	"time"

	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type Poll struct {
	Id          string
	Question    string
	Options     []Option
	CreatedBy   string
	WorkspaceId entity.WorkspaceId
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

type Option struct {
	Id     string
	Text   string
	Votes  int
	Voters []string
}
