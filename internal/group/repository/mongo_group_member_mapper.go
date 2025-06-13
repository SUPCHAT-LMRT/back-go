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
	entityGroupMember *entity.GroupMember,
) (*MongoGroupMember, error) {
	memberObjectId, err := bson.ObjectIDFromHex(entityGroupMember.Id.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert group member id to object id: %w", err)
	}

	userObjectId, err := bson.ObjectIDFromHex(entityGroupMember.UserId.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert user id to object id: %w", err)
	}

	return &MongoGroupMember{
		Id:     memberObjectId,
		UserId: userObjectId,
	}, nil
}

func (m MongoGroupMemberMapper) MapToEntity(mongo *MongoGroupMember) (*entity.GroupMember, error) {
	return &entity.GroupMember{
		Id:     entity.GroupMemberId(mongo.Id.Hex()),
		UserId: user_entity.UserId(mongo.UserId.Hex()),
	}, nil
}
