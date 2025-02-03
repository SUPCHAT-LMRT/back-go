package entity

import (
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ForgotPasswordRequest struct {
	User  *entity.User
	Token uuid.UUID
}
