package repository

import (
	"github.com/supchat-lmrt/back-go/internal/mapper"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ChatDirectMapper struct{}

func NewChatDirectMapper() mapper.Mapper[*MongoChatDirect, *chat_direct_entity.ChatDirect] {
	return &ChatDirectMapper{}
}

func (m ChatDirectMapper) MapToEntity(
	mongo *MongoChatDirect,
) (*chat_direct_entity.ChatDirect, error) {
	reactions := make([]*chat_direct_entity.DirectMessageReaction, len(mongo.Reactions))
	for i, reaction := range mongo.Reactions {
		reactionUsers := make([]entity.UserId, len(reaction.Users))
		for j, user := range reaction.Users {
			reactionUsers[j] = entity.UserId(user.Hex())
		}

		reactions[i] = &chat_direct_entity.DirectMessageReaction{
			Id:       chat_direct_entity.DirectMessageReactionId(reaction.Id.Hex()),
			UserIds:  reactionUsers,
			Reaction: reaction.Reaction,
		}
	}

	attachments := make([]*chat_direct_entity.ChatDirectAttachment, len(mongo.Attachments))
	for i, attachment := range mongo.Attachments {
		attachments[i] = &chat_direct_entity.ChatDirectAttachment{
			Id:       chat_direct_entity.ChatDirectAttachmentId(attachment.Id.Hex()),
			FileName: attachment.FileName,
		}
	}

	return &chat_direct_entity.ChatDirect{
		Id:          chat_direct_entity.ChatDirectId(mongo.Id.Hex()),
		SenderId:    entity.UserId(mongo.SenderId.Hex()),
		User1Id:     entity.UserId(mongo.User1Id.Hex()),
		User2Id:     entity.UserId(mongo.User2Id.Hex()),
		Content:     mongo.Content,
		Reactions:   reactions,
		Attachments: attachments,
		CreatedAt:   mongo.CreatedAt,
		UpdatedAt:   mongo.UpdatedAt,
	}, nil
}

//nolint:revive
func (m ChatDirectMapper) MapFromEntity(
	entity *chat_direct_entity.ChatDirect,
) (*MongoChatDirect, error) {
	chatObjectId, err := bson.ObjectIDFromHex(string(entity.Id))
	if err != nil {
		return nil, err
	}

	senderId, err := bson.ObjectIDFromHex(string(entity.SenderId))
	if err != nil {
		return nil, err
	}

	user1Id, err := bson.ObjectIDFromHex(string(entity.User1Id))
	if err != nil {
		return nil, err
	}

	user2Id, err := bson.ObjectIDFromHex(string(entity.User2Id))
	if err != nil {
		return nil, err
	}

	reactions := make([]*MongoChatDirectReaction, len(entity.Reactions))
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

		reactions[i] = &MongoChatDirectReaction{
			Id:       reactionObjectId,
			Users:    reactionUsers,
			Reaction: reaction.Reaction,
		}
	}

	attachments := make([]*MongoChatDirectAttachment, len(entity.Attachments))
	for i, attachment := range entity.Attachments {
		attachmentObjectId, err := bson.ObjectIDFromHex(string(attachment.Id))
		if err != nil {
			return nil, err
		}

		attachments[i] = &MongoChatDirectAttachment{
			Id:       attachmentObjectId,
			FileName: attachment.FileName,
		}
	}

	return &MongoChatDirect{
		Id:          chatObjectId,
		SenderId:    senderId,
		User1Id:     user1Id,
		User2Id:     user2Id,
		Content:     entity.Content,
		Reactions:   reactions,
		Attachments: attachments,
		UpdatedAt:   entity.UpdatedAt,
		CreatedAt:   entity.CreatedAt,
	}, nil
}
