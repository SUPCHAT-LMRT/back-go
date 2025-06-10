package repository

import (
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoGroupMapper struct{}

func NewMongoGroupMapper() mapper.Mapper[*MongoGroup, *entity.Group] {
	return &MongoGroupMapper{}
}

func (m MongoGroupMapper) MapFromEntity(entityGroup *entity.Group) (*MongoGroup, error) {
	groupObjectId, err := bson.ObjectIDFromHex(entityGroup.Id.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert group id to object id: %w", err)
	}

	ownerObjectId, err := bson.ObjectIDFromHex(entityGroup.OwnerUserId.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert owner id to object id: %w", err)
	}

	return &MongoGroup{
		Id:        groupObjectId,
		Name:      entityGroup.Name,
		OwnerId:   ownerObjectId,
		CreatedAt: entityGroup.CreatedAt,
		UpdatedAt: entityGroup.UpdatedAt,
	}, nil
}

func (m MongoGroupMapper) MapToEntity(mongo *MongoGroup) (*entity.Group, error) {
	return &entity.Group{
		Id:          entity.GroupId(mongo.Id.Hex()),
		Name:        mongo.Name,
		OwnerUserId: user_entity.UserId(mongo.OwnerId.Hex()),
		CreatedAt:   mongo.CreatedAt,
		UpdatedAt:   mongo.UpdatedAt,
	}, nil
}
