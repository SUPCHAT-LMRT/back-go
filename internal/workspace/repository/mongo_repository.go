package repository

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	uberdig "go.uber.org/dig"
	"strings"
)

var (
	databaseName   = "supchat"
	collectionName = "workspaces"
)

type MongoWorkspaceRepositoryDeps struct {
	uberdig.In
	Client                *mongo.Client
	WorkspaceMapper       mapper.Mapper[*MongoWorkspace, *entity.Workspace]
	WorkspaceMemberMapper mapper.Mapper[*MongoWorkspaceMember, *entity.WorkspaceMember]
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

type MongoWorkspaceMember struct {
	Id     bson.ObjectID `bson:"_id"`
	UserId bson.ObjectID `bson:"user_id"`
	Pseudo string        `bson:"pseudo"`
}

func NewMongoWorkspaceRepository(deps MongoWorkspaceRepositoryDeps) WorkspaceRepository {
	return &MongoWorkspaceRepository{deps: deps}
}

func (m MongoWorkspaceRepository) Create(ctx context.Context, workspace *entity.Workspace, ownerMember *entity.WorkspaceMember) error {
	workspace.Id = entity.WorkspaceId(bson.NewObjectID().Hex())
	mongoWorkspace, err := m.deps.WorkspaceMapper.MapFromEntity(workspace)
	if err != nil {
		return err
	}

	session, err := m.deps.Client.Client.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessionContext context.Context) (interface{}, error) {
		_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).InsertOne(sessionContext, mongoWorkspace)
		if err != nil {
			return nil, err
		}

		// Add the owner as a member
		err = m.unsafeAddMember(sessionContext, workspace.Id, ownerMember)
		if err != nil {
			return nil, err
		}

		return nil, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (m MongoWorkspaceRepository) GetById(ctx context.Context, id entity.WorkspaceId) (*entity.Workspace, error) {
	var mongoWorkspace MongoWorkspace

	err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).FindOne(ctx, bson.M{"_id": id}).Decode(&mongoWorkspace)
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

func (m MongoWorkspaceRepository) ListByUserId(ctx context.Context, userId user_entity.UserId) ([]*entity.Workspace, error) {
	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return nil, err
	}

	// Filtre pour trouver les workspaces où l'userId est présent dans members.user_id
	filter := bson.M{"members.user_id": userObjectId}

	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Find(ctx, filter)
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

	// Vérifier les erreurs éventuelles après l'itération
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return workspaces, nil
}

func (m MongoWorkspaceRepository) ListMembers(ctx context.Context, workspaceId entity.WorkspaceId) ([]*entity.WorkspaceMember, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return nil, err
	}

	pipeline := mongo2.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": workspaceObjectId}}},
		{{Key: "$project", Value: bson.M{"_id": 0, "members": 1}}},
		{{Key: "$unwind", Value: "$members"}},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$members"}}},
	}

	// Execute aggregation
	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var members []*entity.WorkspaceMember
	for cursor.Next(ctx) {
		var mongoWorkspaceMember MongoWorkspaceMember
		err = cursor.Decode(&mongoWorkspaceMember)
		if err != nil {
			return nil, err
		}

		member, err := m.deps.WorkspaceMemberMapper.MapToEntity(&mongoWorkspaceMember)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}

	return members, nil
}

func (m MongoWorkspaceRepository) AddMember(ctx context.Context, workspaceId entity.WorkspaceId, member *entity.WorkspaceMember) error {
	member.Id = entity.WorkspaceMemberId(bson.NewObjectID().Hex())
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return err
	}

	// Check if the user is already in the workspace
	workspaceMember, err := m.GetMemberByUserId(ctx, workspaceId, member.UserId)
	if err != nil && !errors.Is(err, WorkspaceMemberNotFoundErr) {
		return err
	}

	if workspaceMember != nil {
		return WorkspaceMemberExistsErr
	}

	mappedMember, err := m.deps.WorkspaceMemberMapper.MapFromEntity(member)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).UpdateOne(ctx, bson.M{"_id": workspaceObjectId}, bson.M{"$addToSet": bson.M{"members": mappedMember}})
	if err != nil {
		return err
	}

	return nil
}

func (m MongoWorkspaceRepository) unsafeAddMember(ctx context.Context, workspaceId entity.WorkspaceId, member *entity.WorkspaceMember) error {
	member.Id = entity.WorkspaceMemberId(bson.NewObjectID().Hex())
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return err
	}

	mappedMember, err := m.deps.WorkspaceMemberMapper.MapFromEntity(member)
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).Collection(collectionName).UpdateOne(ctx, bson.M{"_id": workspaceObjectId}, bson.M{"$addToSet": bson.M{"members": mappedMember}})
	if err != nil {
		return err
	}

	return nil
}

func (m MongoWorkspaceRepository) GetMemberByUserId(ctx context.Context, workspaceId entity.WorkspaceId, userId user_entity.UserId) (*entity.WorkspaceMember, error) {
	// Get the workspace member only
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return nil, err
	}

	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return nil, err
	}

	// MongoDB aggregation pipeline
	pipeline := mongo2.Pipeline{
		{{Key: "$match", Value: bson.M{"_id": workspaceObjectId}}},
		{
			{Key: "$project", Value: bson.M{
				"member": bson.M{
					"$arrayElemAt": []any{
						bson.M{
							"$filter": bson.M{
								"input": "$members",
								"as":    "member",
								"cond":  bson.M{"$eq": []any{"$$member.user_id", userObjectId}},
							},
						},
						0,
					},
				},
			}},
		},
		{{Key: "$replaceRoot", Value: bson.M{"newRoot": "$member"}}},
	}

	cursor, err := m.deps.Client.Client.Database(databaseName).Collection(collectionName).Aggregate(ctx, pipeline)
	if err != nil {
		// Handle (Location40228) PlanExecutor error during aggregation :: caused by :: 'newRoot' expression  must evaluate to an object, but resulting value was: MISSING. Type of resulting value: 'missing'. Input document: {_id: 679b9beb4e1ac8e34475f53c} error
		if strings.Contains(err.Error(), "but resulting value was: MISSING") {
			return nil, WorkspaceMemberNotFoundErr
		}
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoWorkspaceMember MongoWorkspaceMember
	if cursor.Next(ctx) {
		if err = cursor.Decode(&mongoWorkspaceMember); err != nil {
			return nil, err
		}

		member, err := m.deps.WorkspaceMemberMapper.MapToEntity(&mongoWorkspaceMember)
		if err != nil {
			return nil, err
		}

		return member, nil
	}

	return nil, nil
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
