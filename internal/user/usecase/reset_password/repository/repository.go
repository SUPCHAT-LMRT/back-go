package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

var (
	ResetPasswordRequestTtl = 15 * time.Minute
)

var ResetPasswordRequestNotFoundErr = errors.New("validation request not found")

type ResetPasswordRepository interface {
	CreateResetPasswordRequest(ctx context.Context, userId entity.UserId) (*ResetPasswordRequestData, error)
	DeleteResetPasswordRequest(ctx context.Context, validationToken uuid.UUID) (entity.UserId, error)
}

type ResetPasswordRequestData struct {
	UserId entity.UserId
	Token  uuid.UUID
}
