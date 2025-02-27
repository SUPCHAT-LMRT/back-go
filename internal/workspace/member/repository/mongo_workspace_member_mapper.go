package repository

import (
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoWorkspaceMemberMapper struct{}

func NewMongoWorkspaceMemberMapper() mapper.Mapper[*MongoWorkspaceMember, *entity.WorkspaceMember] {
	return &MongoWorkspaceMemberMapper{}
}

func (m MongoWorkspaceMemberMapper) MapFromEntity(entity *entity.WorkspaceMember) (*MongoWorkspaceMember, error) {
	memberObjectId, err := bson.ObjectIDFromHex(entity.Id.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert workspace member id to object id: %w", err)
	}

	userObjectId, err := bson.ObjectIDFromHex(entity.UserId.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert user id to object id: %w", err)
	}

	return &MongoWorkspaceMember{
		Id:     memberObjectId,
		UserId: userObjectId,
		Pseudo: entity.Pseudo,
	}, nil
}

func (m MongoWorkspaceMemberMapper) MapToEntity(databaseWorkspace *MongoWorkspaceMember) (*entity.WorkspaceMember, error) {
	return &entity.WorkspaceMember{
		Id:     entity.WorkspaceMemberId(databaseWorkspace.Id.Hex()),
		UserId: user_entity.UserId(databaseWorkspace.UserId.Hex()),
		Pseudo: databaseWorkspace.Pseudo,
	}, nil
}
