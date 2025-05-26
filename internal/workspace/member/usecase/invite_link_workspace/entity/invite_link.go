package entity

import (
	"time"

	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type InviteLink struct {
	Token       string
	WorkspaceId entity.WorkspaceId
	ExpiresAt   time.Time
}
