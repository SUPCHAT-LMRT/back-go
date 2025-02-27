package entity

import (
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"time"
)

type InviteLink struct {
	Token       string
	WorkspaceId entity.WorkspaceId
	ExpiresAt   time.Time
}
