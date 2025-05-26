package repository

import (
	"context"
	"errors"

	"github.com/google/uuid"
	redis2 "github.com/redis/go-redis/v9"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

var buildResetPasswordRequest = func(token uuid.UUID) string {
	return "reset_password:" + token.String()
}

type RedisResetPasswordRepository struct {
	client *redis.Client
}

func NewRedisResetPasswordRepository(client *redis.Client) ResetPasswordRepository {
	return &RedisResetPasswordRepository{client: client}
}

func (r *RedisResetPasswordRepository) CreateResetPasswordRequest(
	ctx context.Context,
	userId entity.UserId,
) (*ResetPasswordRequestData, error) {
	token := uuid.New()

	err := r.client.Client.Set(ctx, buildResetPasswordRequest(token), userId.String(), ResetPasswordRequestTtl).
		Err()
	if err != nil {
		return nil, err
	}

	return &ResetPasswordRequestData{
		UserId: userId,
		Token:  token,
	}, nil
}

func (r *RedisResetPasswordRepository) DeleteResetPasswordRequest(
	ctx context.Context,
	validationToken uuid.UUID,
) (entity.UserId, error) {
	userId, err := r.client.Client.Get(ctx, buildResetPasswordRequest(validationToken)).Result()
	if err != nil {
		if errors.Is(err, redis2.Nil) {
			return "", ErrResetPasswordRequestNotFound
		}
		return "", err
	}

	err = r.client.Client.Del(ctx, buildResetPasswordRequest(validationToken)).Err()
	if err != nil {
		return "", err
	}

	return entity.UserId(userId), nil
}
