package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

var (
	ForgotPasswordRequestTtl = 5 * time.Minute
)

var ForgotPasswordRequestNotFoundErr = errors.New("validation request not found")

type ForgotPasswordRepository interface {
	CreateForgotPasswordRequest(ctx context.Context, userId entity.UserId) (*ForgotPasswordRequestData, error)
	DeleteForgotPasswordRequest(ctx context.Context, validationToken uuid.UUID) (entity.UserId, error)
}

type ForgotPasswordRequestData struct {
	UserId entity.UserId
	Token  uuid.UUID
}
