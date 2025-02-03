package repository

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

var ValidationRequestNotFoundErr = errors.New("validation request not found")

type ValidationRepository interface {
	CreateValidationRequest(ctx context.Context, userId entity.UserId) (*ValidationRequestData, error)
	DeleteValidationRequest(ctx context.Context, validationToken uuid.UUID) (entity.UserId, error)
}

type ValidationRequestData struct {
	UserId entity.UserId
	Token  uuid.UUID
}
