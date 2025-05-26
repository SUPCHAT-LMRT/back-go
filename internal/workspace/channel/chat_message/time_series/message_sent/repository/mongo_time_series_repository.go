package repository

import (
	"context"
	"errors"
	"time"

	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	databaseName   = "supchat"
	collectionName = "workspace_message_sent_ts"
)

type MongoMessageSentTimeSeriesWorkspaceRepository struct {
	client *mongo.Client
}

func NewMongoMessageSentTimeSeriesWorkspaceRepository(
	client *mongo.Client,
) MessageSentTimeSeriesWorkspaceRepository {
	return &MongoMessageSentTimeSeriesWorkspaceRepository{client: client}
}

func (r MongoMessageSentTimeSeriesWorkspaceRepository) Create(
	ctx context.Context,
	joinedAt time.Time,
	metadata entity.MessageSentMetadata,
) error {
	workspaceObjectId, err := bson.ObjectIDFromHex(string(metadata.WorkspaceId))
	if err != nil {
		return err
	}

	workspaceMemberObjectId, err := bson.ObjectIDFromHex(string(metadata.WorkspaceId))
	if err != nil {
		return err
	}

	channelObjectId, err := bson.ObjectIDFromHex(string(metadata.ChannelId))
	if err != nil {
		return err
	}

	_, err = r.client.Client.Database(databaseName).
		Collection(collectionName).
		InsertOne(ctx, bson.M{
			"metadata": bson.M{
				"workspace_id":     workspaceObjectId,
				"channel_id":       channelObjectId,
				"author_member_id": workspaceMemberObjectId,
			},
			"sent_at": joinedAt,
		})
	if err != nil {
		return err
	}

	return nil
}

func (r MongoMessageSentTimeSeriesWorkspaceRepository) GetMinutelyByWorkspace(
	ctx context.Context,
	workspaceId workspace_entity.WorkspaceId,
	from, to time.Time,
) ([]*entity.MessageSent, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(string(workspaceId))
	if err != nil {
		return nil, err
	}

	pipeline := mongo2.Pipeline{
		bson.D{{Key: "$match", Value: bson.D{
			{Key: "metadata.workspace_id", Value: workspaceObjectId},
		}}},
		bson.D{{Key: "$group", Value: bson.D{
			{Key: "_id", Value: bson.D{
				{Key: "year", Value: bson.D{{Key: "$year", Value: "$sent_at"}}},
				{Key: "month", Value: bson.D{{Key: "$month", Value: "$sent_at"}}},
				{Key: "day", Value: bson.D{{Key: "$dayOfMonth", Value: "$sent_at"}}},
				{Key: "hour", Value: bson.D{{Key: "$hour", Value: "$sent_at"}}},
				{Key: "minute", Value: bson.D{{Key: "$minute", Value: "$sent_at"}}},
			}},
			{Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}},
			{Key: "metadata", Value: bson.D{{Key: "$first", Value: "$metadata"}}},
		}}},
		bson.D{{Key: "$sort", Value: bson.D{
			{Key: "_id.year", Value: 1},
			{Key: "_id.month", Value: 1},
			{Key: "_id.day", Value: 1},
			{Key: "_id.hour", Value: 1},
			{Key: "_id.minute", Value: 1},
		}}},
	}

	cursor, err := r.client.Client.Database(databaseName).
		Collection(collectionName).
		Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var messageSents []*entity.MessageSent
	for cursor.Next(ctx) {
		var result bson.M
		if err = cursor.Decode(&result); err != nil {
			return nil, err
		}

		metadataD, ok := result["metadata"].(bson.D)
		if !ok {
			return nil, errors.New("unable to get metadata")
		}

		metadata := make(bson.M, len(metadataD))
		for _, e := range metadataD {
			metadata[e.Key] = e.Value
		}

		messageSent := entity.MessageSent{
			Metadata: entity.MessageSentMetadata{
				WorkspaceId: workspace_entity.WorkspaceId(
					metadata["workspace_id"].(bson.ObjectID).Hex(),
				),
				ChannelId: channel_entity.ChannelId(
					metadata["channel_id"].(bson.ObjectID).Hex(),
				),
				AuthorMemberId: entity2.WorkspaceMemberId(
					metadata["author_member_id"].(bson.ObjectID).Hex(),
				),
			},
			Count: uint(result["count"].(int32)),
			SentAt: time.Date(
				int(result["_id"].(bson.D)[0].Value.(int32)),
				time.Month(result["_id"].(bson.D)[1].Value.(int32)),
				int(result["_id"].(bson.D)[2].Value.(int32)),
				int(result["_id"].(bson.D)[3].Value.(int32)),
				int(result["_id"].(bson.D)[4].Value.(int32)),
				0,
				0,
				time.UTC,
			),
		}

		messageSents = append(messageSents, &messageSent)
	}

	return messageSents, nil
}
