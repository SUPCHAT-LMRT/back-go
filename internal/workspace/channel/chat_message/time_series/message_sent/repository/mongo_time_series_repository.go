package repository

import (
	"context"
	"errors"
	"github.com/supchat-lmrt/back-go/internal/mongo"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/time_series/message_sent/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongo2 "go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

const (
	databaseName   = "supchat"
	collectionName = "workspace_message_sent_ts"
)

type MongoMessageSentTimeSeriesWorkspaceRepository struct {
	client *mongo.Client
}

func NewMongoMessageSentTimeSeriesWorkspaceRepository(client *mongo.Client) MessageSentTimeSeriesWorkspaceRepository {
	return &MongoMessageSentTimeSeriesWorkspaceRepository{client: client}
}

func (r MongoMessageSentTimeSeriesWorkspaceRepository) Create(ctx context.Context, joinedAt time.Time, metadata entity.MessageSentMetadata) error {
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

	_, err = r.client.Client.Database(databaseName).Collection(collectionName).InsertOne(ctx, bson.M{
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

func (r MongoMessageSentTimeSeriesWorkspaceRepository) GetMinutelyByWorkspace(ctx context.Context, workspaceId workspace_entity.WorkspaceId, from, to time.Time) ([]*entity.MessageSent, error) {
	workspaceObjectId, err := bson.ObjectIDFromHex(string(workspaceId))
	if err != nil {
		return nil, err
	}

	pipeline := mongo2.Pipeline{
		bson.D{{"$match", bson.D{
			{"metadata.workspace_id", workspaceObjectId},
		}}},
		bson.D{{"$group", bson.D{
			{"_id", bson.D{
				{"year", bson.D{{"$year", "$sent_at"}}},
				{"month", bson.D{{"$month", "$sent_at"}}},
				{"day", bson.D{{"$dayOfMonth", "$sent_at"}}},
				{"hour", bson.D{{"$hour", "$sent_at"}}},
				{"minute", bson.D{{"$minute", "$sent_at"}}},
			}},
			{"count", bson.D{{"$sum", 1}}},
			{"metadata", bson.D{{"$first", "$metadata"}}},
		}}},
		bson.D{{"$sort", bson.D{
			{"_id.year", 1},
			{"_id.month", 1},
			{"_id.day", 1},
			{"_id.hour", 1},
			{"_id.minute", 1},
		}}},
	}

	cursor, err := r.client.Client.Database(databaseName).Collection(collectionName).Aggregate(ctx, pipeline)
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
				WorkspaceId:    workspace_entity.WorkspaceId(metadata["workspace_id"].(bson.ObjectID).Hex()),
				ChannelId:      channel_entity.ChannelId(metadata["channel_id"].(bson.ObjectID).Hex()),
				AuthorMemberId: workspace_entity.WorkspaceMemberId(metadata["author_member_id"].(bson.ObjectID).Hex()),
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
