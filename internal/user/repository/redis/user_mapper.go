package redis

import (
	"time"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RedisUserMapper struct{}

func NewRedisUserMapper() mapper.Mapper[map[string]string, *entity.User] {
	return &RedisUserMapper{}
}

func (m RedisUserMapper) MapFromEntity(entityUser *entity.User) (map[string]string, error) {
	userObjectId, err := bson.ObjectIDFromHex(entityUser.Id.String())
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Id":        userObjectId.Hex(),
		"FirstName": entityUser.FirstName,
		"LastName":  entityUser.LastName,
		"Email":     entityUser.Email,
		"CreatedAt": entityUser.CreatedAt.Format(time.RFC3339),
		"UpdatedAt": entityUser.UpdatedAt.Format(time.RFC3339),
	}, nil
}

func (m RedisUserMapper) MapToEntity(databaseUser map[string]string) (*entity.User, error) {
	parsedCreatedAtTime, err := time.Parse(time.RFC3339, databaseUser["CreatedAt"])
	if err != nil {
		return nil, err
	}

	parsedUpdatedAtTime, err := time.Parse(time.RFC3339, databaseUser["UpdatedAt"])
	if err != nil {
		return nil, err
	}

	return &entity.User{
		Id:        entity.UserId(databaseUser["Id"]),
		FirstName: databaseUser["FirstName"],
		LastName:  databaseUser["LastName"],
		Email:     databaseUser["Email"],
		CreatedAt: parsedCreatedAtTime,
		UpdatedAt: parsedUpdatedAtTime,
	}, nil
}
