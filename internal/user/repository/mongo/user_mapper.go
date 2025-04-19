package mongo

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoUserMapper struct{}

func NewMongoUserMapper() mapper.Mapper[*MongoUser, *entity.User] {
	return &MongoUserMapper{}
}

func (m MongoUserMapper) MapFromEntity(entity *entity.User) (*MongoUser, error) {
	userObjectId, err := bson.ObjectIDFromHex(entity.Id.String())
	if err != nil {
		return nil, err
	}

	return &MongoUser{
		Id:         userObjectId,
		FirstName:  entity.FirstName,
		LastName:   entity.LastName,
		Email:      entity.Email,
		OauthEmail: entity.OauthEmail,
		Password:   entity.Password,
		CreatedAt:  entity.CreatedAt,
		UpdatedAt:  entity.UpdatedAt,
	}, nil
}

func (m MongoUserMapper) MapToEntity(databaseUser *MongoUser) (*entity.User, error) {
	return &entity.User{
		Id:         entity.UserId(databaseUser.Id.Hex()),
		FirstName:  databaseUser.FirstName,
		LastName:   databaseUser.LastName,
		Email:      databaseUser.Email,
		OauthEmail: databaseUser.OauthEmail,
		Password:   databaseUser.Password,
		CreatedAt:  databaseUser.CreatedAt,
		UpdatedAt:  databaseUser.UpdatedAt,
	}, nil
}
