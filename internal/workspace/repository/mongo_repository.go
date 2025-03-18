package repository

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	uberdig "go.uber.org/dig"
)

var (
	databaseName          = "supchat"
	collectionName        = "workspaces"
	membersCollectionName = "workspace_members"
)

type MongoWorkspaceRepositoryDeps struct {
	uberdig.In
	Client          *mongo.Client
	WorkspaceMapper mapper.Mapper[*MongoWorkspace, *entity.Workspace]
}

type MongoWorkspaceRepository struct {
	deps MongoWorkspaceRepositoryDeps
}

type MongoWorkspace struct {
	Id      bson.ObjectID `bson:"_id"`
	Name    string        `bson:"name"`
	Type    string        `bson:"type"`
	OwnerId bson.ObjectID `bson:"owner_id"`
}

func NewMongoWorkspaceRepository(deps MongoWorkspaceRepositoryDeps) WorkspaceRepository {
	return &MongoWorkspaceRepository{deps: deps}
}

func (m MongoWorkspaceRepository) Create(ctx context.Context, workspace *entity.Workspace, ownerMember *entity2.WorkspaceMember) error {
	workspace.Id = entity.WorkspaceId(bson.NewObjectID().Hex())
	ownerMember.WorkspaceId = workspace.Id

	mongoWorkspace, err := m.deps.WorkspaceMapper.MapFromEntity(workspace)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).InsertOne(ctx, mongoWorkspace)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoWorkspaceRepository) GetById(ctx context.Context, id entity.WorkspaceId) (*entity.Workspace, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(id.String())
	if err != nil {
		return nil, err
	}

	var mongoWorkspace MongoWorkspace

	err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.M{"_id": workspaceObjectId}).Decode(&mongoWorkspace)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, WorkspaceNotFoundErr
		}
		return nil, err
	}

	workspace, err := m.deps.WorkspaceMapper.MapToEntity(&mongoWorkspace)
	if err != nil {
		return nil, err
	}

	return workspace, nil
}

func (m MongoWorkspaceRepository) ExistsById(ctx context.Context, id entity.WorkspaceId) (bool, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(id.String())
	if err != nil {
		return false, err
	}

	count, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).CountDocuments(ctx, bson.M{"_id": workspaceObjectId})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m MongoWorkspaceRepository) List(ctx context.Context) ([]*entity.Workspace, error) {
	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var workspaces []*entity.Workspace
	for cursor.Next(ctx) {
		var mongoWorkspace MongoWorkspace
		err = cursor.Decode(&mongoWorkspace)
		if err != nil {
			return nil, err
		}

		workspace, err := m.deps.WorkspaceMapper.MapToEntity(&mongoWorkspace)
		if err != nil {
			return nil, err
		}

		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

func (m MongoWorkspaceRepository) ListPublics(ctx context.Context) ([]*entity.Workspace, error) {
	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Find(ctx, bson.M{"type": entity.WorkspaceTypePublic})
	if err != nil {
		return nil, err
	}

	var workspaces []*entity.Workspace
	for cursor.Next(ctx) {
		var mongoWorkspace MongoWorkspace
		err = cursor.Decode(&mongoWorkspace)
		if err != nil {
			return nil, err
		}

		workspace, err := m.deps.WorkspaceMapper.MapToEntity(&mongoWorkspace)
		if err != nil {
			return nil, err
		}

		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

func (m MongoWorkspaceRepository) ListByUserId(ctx context.Context, userId user_entity.UserId) ([]*entity.Workspace, error) {
	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return nil, err
	}

	// Find all workspace IDs where the user is a member
	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(membersCollectionName).Find(ctx, bson.M{"user_id": userObjectId})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var workspaceIds []bson.ObjectID
	for cursor.Next(ctx) {
		var mongoWorkspaceMember struct {
			WorkspaceId bson.ObjectID `bson:"workspace_id"`
		}
		if err = cursor.Decode(&mongoWorkspaceMember); err != nil {
			return nil, err
		}
		workspaceIds = append(workspaceIds, mongoWorkspaceMember.WorkspaceId)
	}

	if len(workspaceIds) == 0 {
		return []*entity.Workspace{}, nil
	}

	// Find all workspaces with the collected workspace IDs
	cursor, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).Find(ctx, bson.M{"_id": bson.M{"$in": workspaceIds}})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var workspaces []*entity.Workspace
	for cursor.Next(ctx) {
		var mongoWorkspace MongoWorkspace
		if err := cursor.Decode(&mongoWorkspace); err != nil {
			return nil, err
		}

		workspace, err := m.deps.WorkspaceMapper.MapToEntity(&mongoWorkspace)
		if err != nil {
			return nil, err
		}

		workspaces = append(workspaces, workspace)
	}

	return workspaces, nil
}

func (m MongoWorkspaceRepository) Update(ctx context.Context, workspace *entity.Workspace) error {
	mongoWorkspace, err := m.deps.WorkspaceMapper.MapFromEntity(workspace)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).UpdateOne(ctx, bson.M{"_id": workspace.Id}, bson.M{"$set": mongoWorkspace})
	if err != nil {
		return err
	}

	return nil
}

func (m MongoWorkspaceRepository) Delete(ctx context.Context, id entity.WorkspaceId) error {
	_, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}
