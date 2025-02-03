package entity

import (
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

var ValidationExpirationTime = 24 * time.Hour

type ValidationRequest struct {
	User  *entity.User
	Token uuid.UUID
}
