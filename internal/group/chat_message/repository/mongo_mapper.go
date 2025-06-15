package repository

import (
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	"github.com/supchat-lmrt/back-go/internal/mapper"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type MongoGroupChatMessageMapper struct{}

func NewMongoGroupChatMessageMapper() mapper.Mapper[*MongoGroupChatMessage, *entity.GroupChatMessage] {
	return &MongoGroupChatMessageMapper{}
}

func (m *MongoGroupChatMessageMapper) MapToEntity(
	model *MongoGroupChatMessage,
) (*entity.GroupChatMessage, error) {
	if model == nil {
		return nil, nil
	}

	messageReactions := make([]*entity.MessageReaction, 0, len(model.Reactions))
	for _, mongoReaction := range model.Reactions {
		userIds := make([]user_entity.UserId, 0, len(mongoReaction.Users))
		for _, userId := range mongoReaction.Users {
			userIds = append(userIds, user_entity.UserId(userId.Hex()))
		}

		messageReactions = append(messageReactions, &entity.MessageReaction{
			Id:       entity.MessageReactionId(mongoReaction.Id.Hex()),
			UserIds:  userIds,
			Reaction: mongoReaction.Reaction,
		})
	}

	attachments := make([]*entity.GroupChatMessageAttachment, len(model.Attachments))
	for i, attachment := range model.Attachments {
		attachments[i] = &entity.GroupChatMessageAttachment{
			Id:       entity.GroupChatAttachmentId(attachment.Id.Hex()),
			FileName: attachment.FileName,
		}
	}

	return &entity.GroupChatMessage{
		Id:          entity.GroupChatMessageId(model.Id.Hex()),
		GroupId:     group_entity.GroupId(model.GroupId.Hex()),
		AuthorId:    user_entity.UserId(model.SenderId.Hex()),
		Content:     model.Content,
		Reactions:   messageReactions,
		Attachments: attachments,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}, nil
}

func (m *MongoGroupChatMessageMapper) MapFromEntity(
	entity *entity.GroupChatMessage,
) (*MongoGroupChatMessage, error) {
	if entity == nil {
		return nil, nil
	}

	msgId, err := bson.ObjectIDFromHex(string(entity.Id))
	if err != nil {
		msgId = bson.NewObjectID()
	}

	groupId, err := bson.ObjectIDFromHex(string(entity.GroupId))
	if err != nil {
		return nil, err
	}

	senderId, err := bson.ObjectIDFromHex(string(entity.AuthorId))
	if err != nil {
		return nil, err
	}

	reactions := make([]*MongoMessageReaction, 0, len(entity.Reactions))
	for _, reaction := range entity.Reactions {
		reactionId, err := bson.ObjectIDFromHex(string(reaction.Id))
		if err != nil {
			reactionId = bson.NewObjectID()
		}

		userObjIds := make([]bson.ObjectID, 0, len(reaction.UserIds))
		for _, userId := range reaction.UserIds {
			userObjId, err := bson.ObjectIDFromHex(string(userId))
			if err != nil {
				return nil, err
			}
			userObjIds = append(userObjIds, userObjId)
		}

		reactions = append(reactions, &MongoMessageReaction{
			Id:       reactionId,
			Users:    userObjIds,
			Reaction: reaction.Reaction,
		})
	}

	attachments := make([]*MongoMessageAttachment, len(entity.Attachments))
	for i, attachment := range entity.Attachments {
		attachmentId, err := bson.ObjectIDFromHex(string(attachment.Id))
		if err != nil {
			attachmentId = bson.NewObjectID()
		}

		attachments[i] = &MongoMessageAttachment{
			Id:       attachmentId,
			FileName: attachment.FileName,
		}
	}

	return &MongoGroupChatMessage{
		Id:          msgId,
		GroupId:     groupId,
		SenderId:    senderId,
		Content:     entity.Content,
		Reactions:   reactions,
		Attachments: attachments,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
	}, nil
}
