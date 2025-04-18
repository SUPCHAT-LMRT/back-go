package repository

import (
	"context"
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/workspace/roles/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "workspace_roles"
)

type MongoRoleRepositoryDeps struct {
	uberdig.In
	Client     *mongo.Client
	RoleMapper mapper.Mapper[*MongoRole, *entity.Role]
}

type MongoRoleRepository struct {
	deps MongoRoleRepositoryDeps
}

type MongoRole struct {
	Id          bson.ObjectID `bson:"_id"`
	Name        string        `bson:"name"`
	WorkspaceId bson.ObjectID `bson:"workspace_id"`
	Permissions uint64        `bson:"permissions"`
	Color       string        `bson:"color"`
}

func NewMongoRoleRepository(deps MongoRoleRepositoryDeps) RoleRepository {
	return &MongoRoleRepository{deps: deps}
}

func (m MongoRoleRepository) Create(ctx context.Context, role *entity.Role) error {
	role.Id = entity.RoleId(bson.NewObjectID().Hex())

	mongoRole, err := m.deps.RoleMapper.MapFromEntity(role)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).InsertOne(ctx, mongoRole)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoRoleRepository) GetById(ctx context.Context, roleId string) (*entity.Role, error) {
	objectId, err := bson.ObjectIDFromHex(roleId)
	if err != nil {
		return nil, fmt.Errorf("invalid role ID: %w", err)
	}

	var mongoRole MongoRole
	err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, bson.M{"_id": objectId}).
		Decode(&mongoRole)
	if err != nil {
		return nil, fmt.Errorf("role not found: %w", err)
	}

	role, err := m.deps.RoleMapper.MapToEntity(&mongoRole)
	if err != nil {
		return nil, fmt.Errorf("error mapping role: %w", err)
	}

	return role, nil
}

func (m MongoRoleRepository) GetList(ctx context.Context, workspaceId string) ([]*entity.Role, error) {
	objectId, err := bson.ObjectIDFromHex(workspaceId)
	if err != nil {
		return nil, fmt.Errorf("invalid workspace ID: %w", err)
	}

	cursor, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Find(ctx, bson.M{"workspace_id": objectId})
	if err != nil {
		return nil, fmt.Errorf("error fetching roles: %w", err)
	}
	defer cursor.Close(ctx)

	var mongoRoles []MongoRole
	if err := cursor.All(ctx, &mongoRoles); err != nil {
		return nil, fmt.Errorf("error decoding roles: %w", err)
	}

	var roles []*entity.Role
	for _, mongoRole := range mongoRoles {
		role, err := m.deps.RoleMapper.MapToEntity(&mongoRole)
		if err != nil {
			return nil, fmt.Errorf("error mapping role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

func (m MongoRoleRepository) Update(ctx context.Context, role *entity.Role) error {
	objectId, err := bson.ObjectIDFromHex(string(role.Id))
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	mongoRole, err := m.deps.RoleMapper.MapFromEntity(role)
	if err != nil {
		return fmt.Errorf("error mapping role: %w", err)
	}

	_, err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.M{"_id": objectId}, bson.M{"$set": mongoRole})
	if err != nil {
		return fmt.Errorf("error updating role: %w", err)
	}

	return nil
}

func (m MongoRoleRepository) Delete(ctx context.Context, roleId string) error {
	objectId, err := bson.ObjectIDFromHex(roleId)
	if err != nil {
		return fmt.Errorf("invalid role ID: %w", err)
	}

	_, err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		DeleteOne(ctx, bson.M{"_id": objectId})
	if err != nil {
		return fmt.Errorf("error deleting role: %w", err)
	}

	return nil
}
