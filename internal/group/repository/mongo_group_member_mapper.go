package repository

import (
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoGroupMemberMapper struct{}

func NewMongoGroupMemberMapper() mapper.Mapper[*MongoGroupMember, *entity.GroupMember] {
	return &MongoGroupMemberMapper{}
}

func (m MongoGroupMemberMapper) MapFromEntity(
	entity *entity.GroupMember,
) (*MongoGroupMember, error) {
	memberObjectId, err := bson.ObjectIDFromHex(entity.Id.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert group member id to object id: %w", err)
	}

	userObjectId, err := bson.ObjectIDFromHex(entity.UserId.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert user id to object id: %w", err)
	}

	groupObjectId, err := bson.ObjectIDFromHex(entity.GroupId.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert group id to object id: %w", err)
	}

	return &MongoGroupMember{
		Id:      memberObjectId,
		UserId:  userObjectId,
		GroupId: groupObjectId,
	}, nil
}

func (m MongoGroupMemberMapper) MapToEntity(mongo *MongoGroupMember) (*entity.GroupMember, error) {
	return &entity.GroupMember{
		Id:      entity.GroupMemberId(mongo.Id.Hex()),
		UserId:  user_entity.UserId(mongo.UserId.Hex()),
		GroupId: entity.GroupId(mongo.GroupId.Hex()),
	}, nil
}
