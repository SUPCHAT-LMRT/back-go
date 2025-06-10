package repository

import (
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoWorkspaceMapper struct{}

func NewMongoWorkspaceMapper() mapper.Mapper[*MongoWorkspace, *entity.Workspace] {
	return &MongoWorkspaceMapper{}
}

func (m MongoWorkspaceMapper) MapFromEntity(
	entityWorkspace *entity.Workspace,
) (*MongoWorkspace, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(entityWorkspace.Id.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert workspace id to object id: %w", err)
	}

	ownerObjectId, err := bson.ObjectIDFromHex(entityWorkspace.OwnerId.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert owner id to object id: %w", err)
	}

	return &MongoWorkspace{
		Id:      workspaceObjectId,
		Name:    entityWorkspace.Name,
		Topic:   entityWorkspace.Topic,
		Type:    string(entityWorkspace.Type),
		OwnerId: ownerObjectId,
	}, nil
}

func (m MongoWorkspaceMapper) MapToEntity(
	databaseWorkspace *MongoWorkspace,
) (*entity.Workspace, error) {
	return &entity.Workspace{
		Id:      entity.WorkspaceId(databaseWorkspace.Id.Hex()),
		Name:    databaseWorkspace.Name,
		Topic:   databaseWorkspace.Topic,
		Type:    entity.WorkspaceType(databaseWorkspace.Type),
		OwnerId: user_entity.UserId(databaseWorkspace.OwnerId.Hex()),
	}, nil
}
