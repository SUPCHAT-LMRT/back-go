package repository

import (
	"context"
	"github.com/google/uuid"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/redis"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/validation/entity"
	uberdig "go.uber.org/dig"
)

var (
	buildValidationRequestRedisKey = func(token uuid.UUID) string {
		return "validation_request:" + token.String()
	}
)

type RedisValidationRepositoryDeps struct {
	uberdig.In
	Client *redis.Client
	Mapper mapper.Mapper[map[string]string, *ValidationRequestData]
}

type RedisValidationRepository struct {
	deps RedisValidationRepositoryDeps
}

func NewRedisValidationRepository(deps RedisValidationRepositoryDeps) ValidationRepository {
	return &RedisValidationRepository{deps: deps}
}

func (m RedisValidationRepository) CreateValidationRequest(ctx context.Context, userId user_entity.UserId) (*ValidationRequestData, error) {
	token := uuid.New()

	err := m.deps.Client.Client.Set(ctx, buildValidationRequestRedisKey(token), userId.String(), entity.ValidationExpirationTime).Err()
	if err != nil {
		return nil, err
	}

	return &ValidationRequestData{
		UserId: userId,
		Token:  token,
	}, nil
}

func (m RedisValidationRepository) DeleteValidationRequest(ctx context.Context, validationToken uuid.UUID) (user_entity.UserId, error) {
	userId, err := m.deps.Client.Client.Get(ctx, buildValidationRequestRedisKey(validationToken)).Result()
	if err != nil {
		return "", err
	}

	err = m.deps.Client.Client.Del(ctx, buildValidationRequestRedisKey(validationToken)).Err()
	if err != nil {
		return "", err
	}

	return user_entity.UserId(userId), nil
}
