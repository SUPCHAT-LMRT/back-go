package repository

import (
	"context"
	"errors"
	"time"

	"github.com/supchat-lmrt/back-go/internal/bots/poll/entity"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
)

var (
	databaseName   = "supchat"
	collectionName = "polls"
)

type MongoPollRepositoryDeps struct {
	uberdig.In
	Client     *mongo.Client
	PollMapper *MongoPollMapper
}

type MongoPoll struct {
	Id          string              `bson:"_id"`
	Question    string              `bson:"question"`
	Options     []Option            `bson:"options"`
	CreatedBy   string              `bson:"created_by"`
	WorkspaceId entity2.WorkspaceId `bson:"workspace_id"`
	CreatedAt   time.Time           `bson:"created_at"`
	ExpiresAt   time.Time           `bson:"expires_at"`
}

type Option struct {
	Id     string   `bson:"id"`
	Text   string   `bson:"text"`
	Votes  int      `bson:"votes"`
	Voters []string `bson:"voters"`
}

type MongoPollRepository struct {
	deps MongoPollRepositoryDeps
}

func NewMongoPollRepository(deps MongoPollRepositoryDeps) PollRepository {
	return &MongoPollRepository{deps: deps}
}

func (r *MongoPollRepository) Create(ctx context.Context, poll *entity.Poll) error {
	mongoPoll, err := r.deps.PollMapper.MapFromEntity(poll)
	if err != nil {
		return err
	}

	_, err = r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		InsertOne(ctx, mongoPoll)
	return err
}

func (r *MongoPollRepository) GetById(ctx context.Context, pollId string) (*entity.Poll, error) {
	var mongoPoll MongoPoll
	err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		FindOne(ctx, bson.M{"_id": pollId}).
		Decode(&mongoPoll)
	if err != nil {
		return nil, errors.New("poll not found")
	}

	return r.deps.PollMapper.MapToEntity(&mongoPoll)
}

func (r *MongoPollRepository) GetAllByWorkspace(
	ctx context.Context,
	workspaceId entity2.WorkspaceId,
) ([]*entity.Poll, error) {
	cursor, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		Find(ctx, bson.M{"workspace_id": workspaceId}) // Filtrer par workspace
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var mongoPolls []MongoPoll
	if err := cursor.All(ctx, &mongoPolls); err != nil {
		return nil, err
	}

	var polls []*entity.Poll
	for _, mongoPoll := range mongoPolls {
		poll, err := r.deps.PollMapper.MapToEntity(&mongoPoll)
		if err != nil {
			return nil, err
		}
		polls = append(polls, poll)
	}

	return polls, nil
}

func (r *MongoPollRepository) Delete(ctx context.Context, pollId string) error {
	_, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		DeleteOne(ctx, bson.M{"_id": pollId})
	if err != nil {
		return errors.New("poll not found")
	}
	return nil
}

func (r *MongoPollRepository) Vote(ctx context.Context, poll *entity.Poll) error {
	mongoPoll, err := r.deps.PollMapper.MapFromEntity(poll)
	if err != nil {
		return err
	}

	result, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, bson.M{"_id": poll.Id}, bson.M{"$set": mongoPoll})
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("poll not found")
	}

	return nil
}

func (r *MongoPollRepository) IncrementVote(
	ctx context.Context,
	pollId string,
	optionId string,
) error {
	filter := bson.M{"_id": pollId, "options.id": optionId}
	update := bson.M{"$inc": bson.M{"options.$.votes": 1}}

	result, err := r.deps.Client.Client.Database(databaseName).
		Collection(collectionName).
		UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return errors.New("poll or option not found")
	}

	return nil
}
