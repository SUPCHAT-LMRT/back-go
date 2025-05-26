package repository

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChannelMessageMapper struct{}

func NewChannelMessageMapper() mapper.Mapper[*MongoChannelMessage, *entity.ChannelMessage] {
	return &ChannelMessageMapper{}
}

func (m ChannelMessageMapper) MapToEntity(
	mongo *MongoChannelMessage,
) (*entity.ChannelMessage, error) {
	reactions := make([]*entity.ChannelMessageReaction, len(mongo.Reactions))
	for i, reaction := range mongo.Reactions {
		reactionUsers := make([]user_entity.UserId, len(reaction.Users))
		for j, user := range reaction.Users {
			reactionUsers[j] = user_entity.UserId(user.Hex())
		}

		reactions[i] = &entity.ChannelMessageReaction{
			Id:        entity.ChannelMessageReactionId(reaction.Id.Hex()),
			MessageId: entity.ChannelMessageId(mongo.Id.Hex()),
			UserIds:   reactionUsers,
			Reaction:  reaction.Reaction,
		}
	}

	return &entity.ChannelMessage{
		Id:        entity.ChannelMessageId(mongo.Id.Hex()),
		ChannelId: channel_entity.ChannelId(mongo.ChannelId.Hex()),
		AuthorId:  user_entity.UserId(mongo.AuthorId.Hex()),
		Content:   mongo.Content,
		CreatedAt: mongo.CreatedAt,
		UpdatedAt: mongo.UpdatedAt,
		Reactions: reactions,
	}, nil
}

func (m ChannelMessageMapper) MapFromEntity(
	entity *entity.ChannelMessage,
) (*MongoChannelMessage, error) {
	messageObjectId, err := bson.ObjectIDFromHex(string(entity.Id))
	if err != nil {
		return nil, err
	}
	channelObjectId, err := bson.ObjectIDFromHex(string(entity.ChannelId))
	if err != nil {
		return nil, err
	}
	authorObjectId, err := bson.ObjectIDFromHex(string(entity.AuthorId))
	if err != nil {
		return nil, err
	}

	reactions := make([]*MongoChannelMessageReaction, len(entity.Reactions))
	for i, reaction := range entity.Reactions {
		reactionUsers := make([]bson.ObjectID, len(reaction.UserIds))
		for j, user := range reaction.UserIds {
			userObjectId, err := bson.ObjectIDFromHex(string(user))
			if err != nil {
				return nil, err
			}

			reactionUsers[j] = userObjectId
		}

		reactionObjectId, err := bson.ObjectIDFromHex(string(reaction.Id))
		if err != nil {
			return nil, err
		}

		reactions[i] = &MongoChannelMessageReaction{
			Id:       reactionObjectId,
			Users:    reactionUsers,
			Reaction: reaction.Reaction,
		}
	}

	return &MongoChannelMessage{
		Id:        messageObjectId,
		ChannelId: channelObjectId,
		AuthorId:  authorObjectId,
		Content:   entity.Content,
		CreatedAt: entity.CreatedAt,
		Reactions: reactions,
	}, nil
}
