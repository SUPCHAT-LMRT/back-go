package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

var buildForgotPasswordRedisKey = func(token uuid.UUID) string {
	return "forgot_password:" + token.String()
}

type RedisForgotPasswordRepository struct {
	client *redis.Client
}

func NewRedisForgotPasswordRepository(client *redis.Client) ForgotPasswordRepository {
	return &RedisForgotPasswordRepository{client: client}
}

func (r *RedisForgotPasswordRepository) CreateForgotPasswordRequest(
	ctx context.Context,
	userId entity.UserId,
) (*ForgotPasswordRequestData, error) {
	token := uuid.New()

	err := r.client.Client.Set(ctx, buildForgotPasswordRedisKey(token), userId.String(), ForgotPasswordRequestTtl).
		Err()
	if err != nil {
		return nil, err
	}

	return &ForgotPasswordRequestData{
		UserId: userId,
		Token:  token,
	}, nil
}

func (r *RedisForgotPasswordRepository) DeleteForgotPasswordRequest(
	ctx context.Context,
	validationToken uuid.UUID,
) (entity.UserId, error) {
	userId, err := r.client.Client.Get(ctx, buildForgotPasswordRedisKey(validationToken)).Result()
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return "", ErrForgotPasswordRequestNotFound
		}
		return "", err
	}

	err = r.client.Client.Del(ctx, buildForgotPasswordRedisKey(validationToken)).Err()
	if err != nil {
		return "", err
	}

	return entity.UserId(userId), nil
}
