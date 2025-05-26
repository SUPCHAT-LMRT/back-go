package repository

import (
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoRoleMapper struct{}

func NewMongoRoleMapper() mapper.Mapper[*MongoRole, *entity.Role] {
	return &MongoRoleMapper{}
}

func (m MongoRoleMapper) MapFromEntity(role *entity.Role) (*MongoRole, error) {
	roleObjectId, err := bson.ObjectIDFromHex(role.Id.String())
	if err != nil {
		return nil, fmt.Errorf("unable to convert workspace id to object id: %w", err)
	}

	workspaceID, err := bson.ObjectIDFromHex(string(role.WorkspaceId))
	if err != nil {
		return nil, fmt.Errorf("unable to convert workspace id to object id: %w", err)
	}

	return &MongoRole{
		Id:          roleObjectId,
		Name:        role.Name,
		WorkspaceId: workspaceID,
		Permissions: role.Permissions,
		Color:       role.Color,
	}, nil
}

func (m MongoRoleMapper) MapToEntity(mongo *MongoRole) (*entity.Role, error) {
	return &entity.Role{
		Id:          entity.RoleId(mongo.Id.Hex()),
		Name:        mongo.Name,
		WorkspaceId: workspace_entity.WorkspaceId(mongo.WorkspaceId.Hex()),
		Permissions: mongo.Permissions,
		Color:       mongo.Color,
	}, nil
}
