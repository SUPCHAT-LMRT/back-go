package redis

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type RedisUserMapper struct{}

func NewRedisUserMapper() mapper.Mapper[map[string]string, *entity.User] {
	return &RedisUserMapper{}
}

func (m RedisUserMapper) MapFromEntity(entity *entity.User) (map[string]string, error) {
	userObjectId, err := bson.ObjectIDFromHex(entity.Id.String())
	if err != nil {
		return nil, err
	}

	return map[string]string{
		"Id":        userObjectId.Hex(),
		"FirstName": entity.FirstName,
		"LastName":  entity.LastName,
		"Email":     entity.Email,
		"CreatedAt": entity.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (m RedisUserMapper) MapToEntity(databaseUser map[string]string) (*entity.User, error) {
	parsedTime, err := time.Parse(time.RFC3339, databaseUser["CreatedAt"])
	if err != nil {
		return nil, err
	}

	return &entity.User{
		Id:        entity.UserId(databaseUser["Id"]),
		FirstName: databaseUser["FirstName"],
		LastName:  databaseUser["LastName"],
		Email:     databaseUser["Email"],
		CreatedAt: parsedTime,
	}, nil
}
