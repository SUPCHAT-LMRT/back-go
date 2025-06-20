package repository

import (
	"context"
	"errors"
	"time"

	"github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "groups"
)

type MongoGroupRepositoryDeps struct {
	uberdig.In
	Client            *mongo.Client
	GroupMapper       mapper.Mapper[*MongoGroup, *entity.Group]
	GroupMemberMapper mapper.Mapper[*MongoGroupMember, *entity.GroupMember]
}

type MongoGroupRepository struct {
	deps MongoGroupRepositoryDeps
}

type MongoGroup struct {
	Id        bson.ObjectID `bson:"_id"`
	Name      string        `bson:"name"`
	OwnerId   bson.ObjectID `bson:"owner_id"`
	CreatedAt time.Time     `bson:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at"`
}

type MongoGroupMember struct {
	Id     bson.ObjectID `bson:"_id"`
	UserId bson.ObjectID `bson:"user_id"`
}

func NewMongoGroupRepository(deps MongoGroupRepositoryDeps) GroupRepository {
	return &MongoGroupRepository{deps: deps}
}

func (r MongoGroupRepository) Create(
	ctx context.Context,
	group *entity.Group,
	ownerUserId user_entity.UserId,
) error {
	group.Id = entity.GroupId(bson.NewObjectID().Hex())
	group.CreatedAt = time.Now()
	group.UpdatedAt = group.CreatedAt
	group.OwnerMemberId = entity.GroupMemberId(bson.NewObjectID().Hex())

	mongoGroup, err := r.deps.GroupMapper.MapFromEntity(group)
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		InsertOne(ctx, mongoGroup)
	if err != nil {
		return err
	}

	// Add the owner as a member
	err = r.unsafeAddMember(ctx, group.Id, &entity.GroupMember{
		Id:     group.OwnerMemberId,
		UserId: ownerUserId,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r MongoGroupRepository) DeleteGroup(ctx context.Context, groupId entity.GroupId) error {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		DeleteOne(ctx, bson.M{"_id": groupObjectId})
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return ErrGroupNotFound
		}
		return err
	}

	return nil
}

func (r MongoGroupRepository) GetGroup(
	ctx context.Context,
	groupId entity.GroupId,
) (*entity.Group, error) {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return nil, err
	}

	var mongoGroup MongoGroup
	err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, bson.M{"_id": groupObjectId}).
		Decode(&mongoGroup)
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	group, err := r.deps.GroupMapper.MapToEntity(&mongoGroup)
	if err != nil {
		return nil, err
	}

	return group, nil
}

func (r MongoGroupRepository) GetMemberByUserId(
	ctx context.Context,
	groupId entity.GroupId,
	userId user_entity.UserId,
) (*entity.GroupMember, error) {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return nil, err
	}
	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return nil, err
	}

	var result struct {
		Members []MongoGroupMember `bson:"members"`
	}
	err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, bson.M{"_id": groupObjectId}).
		Decode(&result)

	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return nil, ErrGroupNotFound
		}
		return nil, err
	}

	for _, member := range result.Members {
		if member.UserId == userObjectId {
			return r.deps.GroupMemberMapper.MapToEntity(&member)
		}
	}

	return nil, errors.New("member not found")
}

func (r MongoGroupRepository) ListRecentGroups(ctx context.Context, userId user_entity.UserId) ([]*entity.Group, error) {
	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return nil, err
	}

	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})
	cursor, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Find(ctx, bson.M{"members.user_id": userObjectId}, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoGroups []*MongoGroup
	err = cursor.All(ctx, &mongoGroups)
	if err != nil {
		return nil, err
	}

	groups := make([]*entity.Group, len(mongoGroups))
	for i, mongoGroup := range mongoGroups {
		group, err := r.deps.GroupMapper.MapToEntity(mongoGroup)
		if err != nil {
			return nil, err
		}
		groups[i] = group
	}

	return groups, nil
}

func (r MongoGroupRepository) AddMember(
	ctx context.Context,
	groupId entity.GroupId,
	inviteeUserId user_entity.UserId,
) error {
	isMember, err := r.isMember(ctx, groupId, inviteeUserId)
	if err != nil {
		return err
	}

	if isMember {
		return ErrMemberAlreadyInGroup
	}

	return r.unsafeAddMember(ctx, groupId, &entity.GroupMember{
		Id:     entity.GroupMemberId(bson.NewObjectID().Hex()),
		UserId: inviteeUserId,
	})
}

func (r MongoGroupRepository) RemoveMember(
	ctx context.Context,
	groupId entity.GroupId,
	memberId entity.GroupMemberId,
) error {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return err
	}
	memberObjectId, err := bson.ObjectIDFromHex(memberId.String())
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx,
			bson.M{"_id": groupObjectId},
			bson.M{"$pull": bson.M{"members": bson.M{"_id": memberObjectId}}})
	if err != nil {
		if errors.Is(err, mongo2.ErrNoDocuments) {
			return ErrGroupNotFound
		}
		return err
	}

	return nil
}

func (r MongoGroupRepository) Exists(ctx context.Context, groupId entity.GroupId) (bool, error) {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return false, err
	}

	count, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, bson.M{"_id": groupObjectId})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r MongoGroupRepository) UpdateGroupName(
	ctx context.Context,
	groupId entity.GroupId,
	name string,
) error {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.M{"_id": groupObjectId}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		return err
	}

	return nil
}

func (r MongoGroupRepository) ListMembers(
	ctx context.Context,
	groupId entity.GroupId,
) ([]*entity.GroupMember, error) {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return nil, err
	}

	pipeline := mongo2.Pipeline{
		{{Key: "$match", Value: bson.D{{Key: "_id", Value: groupObjectId}}}},
		{{Key: "$unwind", Value: "$members"}},
		{{Key: "$replaceRoot", Value: bson.D{{Key: "newRoot", Value: "$members"}}}},
	}

	cursor, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var members []MongoGroupMember
	if err = cursor.All(ctx, &members); err != nil {
		return nil, err
	}

	groupMembers := make([]*entity.GroupMember, len(members))
	for i, member := range members {
		groupMember, err := r.deps.GroupMemberMapper.MapToEntity(&member)
		if err != nil {
			return nil, err
		}

		groupMembers[i] = groupMember
	}

	return groupMembers, nil
}

func (r MongoGroupRepository) isMember(
	ctx context.Context,
	groupId entity.GroupId,
	userId user_entity.UserId,
) (bool, error) {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return false, err
	}
	userObjectId, err := bson.ObjectIDFromHex(userId.String())
	if err != nil {
		return false, err
	}

	count, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, bson.M{"_id": groupObjectId, "members.user_id": userObjectId})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

// unsafeAddMember adds a member to a group without any checks. (e.g. if the user is already a member)
func (r MongoGroupRepository) unsafeAddMember(
	ctx context.Context,
	groupId entity.GroupId,
	member *entity.GroupMember,
) error {
	workspaceObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return err
	}

	mappedMember, err := r.deps.GroupMemberMapper.MapFromEntity(member)
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.M{"_id": workspaceObjectId}, bson.M{"$addToSet": bson.M{"members": mappedMember}})
	if err != nil {
		return err
	}

	return nil
}

func (r MongoGroupRepository) TransferOwnership(
	ctx context.Context,
	groupId entity.GroupId,
	newOwnerMemberId entity.GroupMemberId,
) error {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return err
	}
	newOwnerObjectId, err := bson.ObjectIDFromHex(newOwnerMemberId.String())
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.M{"_id": groupObjectId}, bson.M{"$set": bson.M{"owner_id": newOwnerObjectId}})
	if err != nil {
		return err
	}

	return nil
}

func (r MongoGroupRepository) IsMember(
	ctx context.Context,
	groupId entity.GroupId,
	memberId entity.GroupMemberId,
) (bool, error) {
	groupObjectId, err := bson.ObjectIDFromHex(groupId.String())
	if err != nil {
		return false, err
	}
	memberObjectId, err := bson.ObjectIDFromHex(memberId.String())
	if err != nil {
		return false, err
	}

	count, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		CountDocuments(ctx, bson.M{"_id": groupObjectId, "members._id": memberObjectId})
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
