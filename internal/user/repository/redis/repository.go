package redis

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/redis"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/repository"
	uberdig "go.uber.org/dig"
	"time"
)

var (
	userRedisExpiration = 5 * time.Minute

	buildUserRedisKey = func(userId entity.UserId) string {
		return "user:" + userId.String()
	}
	buildUserEmailBindingRedisKey = func(userEmail string) string {
		return "user_email:" + userEmail
	}
)

type RedisUserRepositoryDeps struct {
	uberdig.In
	Client     *redis.Client
	UserMapper mapper.Mapper[map[string]string, *entity.User]
}

type RedisUserRepository struct {
	deps RedisUserRepositoryDeps
}

func NewRedisUserRepository(deps RedisUserRepositoryDeps) repository.UserRepository {
	return &RedisUserRepository{deps: deps}
}

func (r RedisUserRepository) Create(ctx context.Context, user *entity.User) error {
	redisEntity, err := r.deps.UserMapper.MapFromEntity(user)
	if err != nil {
		return err
	}

	err = r.deps.Client.Client.HSet(ctx, buildUserRedisKey(user.Id), redisEntity).Err()
	if err != nil {
		return err
	}

	err = r.deps.Client.Client.Expire(ctx, buildUserRedisKey(user.Id), userRedisExpiration).Err()
	if err != nil {
		return err
	}

	// Create bindings
	err = r.deps.Client.Client.Set(ctx, buildUserEmailBindingRedisKey(user.Email), user.Id.String(), userRedisExpiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisUserRepository) GetById(ctx context.Context, userId entity.UserId) (user *entity.User, err error) {
	result, err := r.deps.Client.Client.HGetAll(ctx, buildUserRedisKey(userId)).Result()
	if err != nil {
		return nil, err
	}

	user, err = r.deps.UserMapper.MapToEntity(result)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r RedisUserRepository) GetByEmail(ctx context.Context, userEmail string, options ...repository.GetUserOptionFunc) (user *entity.User, err error) {
	userIdStr, err := r.deps.Client.Client.Get(ctx, buildUserEmailBindingRedisKey(userEmail)).Result()
	if err != nil {
		return nil, err
	}

	return r.GetById(ctx, entity.UserId(userIdStr))
}

func (r RedisUserRepository) List(ctx context.Context) ([]*entity.User, error) {
	const batchSize = 100
	var cursor uint64

	users := make([]*entity.User, 0)

	for {
		keys, nextCursor, err := r.deps.Client.Client.Scan(ctx, cursor, "user:*", batchSize).Result()
		if err != nil {
			return nil, err
		}

		for _, key := range keys {
			user, err := r.GetById(ctx, entity.UserId(key[len("user:"):]))
			if err != nil {
				return nil, err
			}
			users = append(users, user)
		}

		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return users, nil
}

func (r RedisUserRepository) Update(ctx context.Context, user *entity.User) error {
	return r.Create(ctx, user)
}

func (r RedisUserRepository) Delete(ctx context.Context, userId entity.UserId) error {
	// Remove bindings
	user, err := r.GetById(ctx, userId)
	if err != nil {
		return err
	}

	err = r.deps.Client.Client.Del(ctx, buildUserEmailBindingRedisKey(user.Email)).Err()
	if err != nil {
		return err
	}

	err = r.deps.Client.Client.Del(ctx, buildUserRedisKey(userId)).Err()
	if err != nil {
		return err
	}

	return nil
}
