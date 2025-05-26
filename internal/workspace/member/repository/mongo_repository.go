package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "workspace_members"
)

type MongoWorkspaceMemberRepositoryDeps struct {
	uberdig.In
	Client                *mongo.Client
	WorkspaceMemberMapper mapper.Mapper[*MongoWorkspaceMember, *entity2.WorkspaceMember]
}

type MongoWorkspaceMemberRepository struct {
	deps MongoWorkspaceMemberRepositoryDeps
}

type MongoWorkspaceMember struct {
	Id          bson.ObjectID `bson:"_id"`
	WorkspaceId bson.ObjectID `bson:"workspace_id"`
	UserId      bson.ObjectID `bson:"user_id"`
}

func NewMongoWorkspaceMemberRepository(
	deps MongoWorkspaceMemberRepositoryDeps,
) WorkspaceMemberRepository {
	return &MongoWorkspaceMemberRepository{deps: deps}
}

//nolint:revive
func (m MongoWorkspaceMemberRepository) ListMembers(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	limit, page int,
) (totalMembers uint, members []*entity2.WorkspaceMember, err error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return 0, nil, err
	}

	// Define query options for pagination
	opts := options.Find().SetLimit(int64(limit)).SetSkip(int64((page - 1) * limit))

	cursor, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Find(ctx, bson.M{"workspace_id": workspaceObjectId}, opts)
	if err != nil {
		return 0, nil, err
	}
	defer cursor.Close(ctx)

	var total int64
	total, err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, bson.M{"workspace_id": workspaceObjectId})
	if err != nil {
		return 0, nil, err
	}

	totalMembers = uint(total)

	var mongoWorkspaceMembers []*MongoWorkspaceMember
	for cursor.Next(ctx) {
		var mongoWorkspaceMember MongoWorkspaceMember
		if err = cursor.Decode(&mongoWorkspaceMember); err != nil {
			return 0, nil, err
		}
		mongoWorkspaceMembers = append(mongoWorkspaceMembers, &mongoWorkspaceMember)
	}

	workspaceMembers := make([]*entity2.WorkspaceMember, len(mongoWorkspaceMembers))
	for i, mongoWorkspaceMember := range mongoWorkspaceMembers {
		workspaceMember, err := m.deps.WorkspaceMemberMapper.MapToEntity(mongoWorkspaceMember)
		if err != nil {
			return 0, nil, err
		}
		workspaceMembers[i] = workspaceMember
	}

	return totalMembers, workspaceMembers, nil
}

func (m MongoWorkspaceMemberRepository) AddMember(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	member *entity2.WorkspaceMember,
) error {
	member.Id = entity2.WorkspaceMemberId(bson.NewObjectID().Hex())
	member.WorkspaceId = workspaceId

	// Check if the user is already in the workspace
	workspaceMember, err := m.GetMemberByUserId(ctx, workspaceId, member.UserId)
	if err != nil && !errors.Is(err, ErrWorkspaceMemberNotFound) {
		return err
	}

	if workspaceMember != nil {
		return ErrWorkspaceMemberExists
	}

	err = m.unsafeAddMember(ctx, member)
	if err != nil {
		return fmt.Errorf("failed to add member: %w", err)
	}

	return nil
}

func (m MongoWorkspaceMemberRepository) unsafeAddMember(
	ctx context.Context,
	member *entity2.WorkspaceMember,
) error {
	member.Id = entity2.WorkspaceMemberId(bson.NewObjectID().Hex())
	mappedMember, err := m.deps.WorkspaceMemberMapper.MapFromEntity(member)
	if err != nil {
		return err
	}

	mappedMember.WorkspaceId, err = bson.ObjectIDFromHex(member.WorkspaceId.String())
	if err != nil {
		return err
	}

	_, err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		InsertOne(ctx, mappedMember)
	if err != nil {
		return err
	}

	return nil
}

func (m MongoWorkspaceMemberRepository) CountMembers(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
) (uint, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return 0, err
	}

	count, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, bson.M{"workspace_id": workspaceObjectId})
	if err != nil {
		return 0, err
	}

	return uint(count), nil
}

func (m MongoWorkspaceMemberRepository) GetMemberByUserId(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	userId user_entity.UserId,
) (*entity2.WorkspaceMember, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return nil, err
	}

	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return nil, err
	}

	var mongoWorkspaceMember MongoWorkspaceMember
	err = m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, bson.M{"workspace_id": workspaceObjectId, "user_id": userObjectId}).
		Decode(&mongoWorkspaceMember)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, ErrWorkspaceMemberNotFound
		}
		return nil, err
	}

	return m.deps.WorkspaceMemberMapper.MapToEntity(&mongoWorkspaceMember)
}

func (m MongoWorkspaceMemberRepository) IsMemberExists(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	memberId entity2.WorkspaceMemberId,
) (bool, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return false, err
	}

	memberObjectId, err := bson.ObjectIDFromHex(memberId.String())
	if err != nil {
		return false, err
	}

	count, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, bson.M{
			"workspace_id": workspaceObjectId,
			"_id":          memberObjectId,
		})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m MongoWorkspaceMemberRepository) IsMemberByUserIdExists(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	userId user_entity.UserId,
) (bool, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return false, err
	}

	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return false, err
	}

	count, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, bson.M{
			"workspace_id": workspaceObjectId,
			"user_id":      userObjectId,
		})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (m MongoWorkspaceMemberRepository) RemoveMember(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	memberId entity2.WorkspaceMemberId,
) error {
	workspaceObjectId, err := bson.ObjectIDFromHex(workspaceId.String())
	if err != nil {
		return err
	}

	memberObjectId, err := bson.ObjectIDFromHex(memberId.String())
	if err != nil {
		return err
	}

	result, err := m.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		DeleteOne(ctx, bson.M{
			"workspace_id": workspaceObjectId,
			"_id":          memberObjectId,
		})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrWorkspaceMemberNotFound
	}

	return nil
}
