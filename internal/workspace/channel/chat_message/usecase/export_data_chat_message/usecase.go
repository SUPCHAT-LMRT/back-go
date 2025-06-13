package export_data_chat_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
	entity2 "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	"time"
)

type ExportDataChatMessageUseCase struct {
	repository repository.ChannelMessageRepository
}

func NewExportDataChatMessageUseCase(repo repository.ChannelMessageRepository) *ExportDataChatMessageUseCase {
	return &ExportDataChatMessageUseCase{repository: repo}
}

func (uc *ExportDataChatMessageUseCase) Execute(ctx context.Context, userId entity.UserId) ([]ExportableChatMessageData, error) {
	messages, err := uc.repository.ListAllMessagesByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	exported := make([]ExportableChatMessageData, len(messages))
	for i, msg := range messages {
		exported[i] = ExportableChatMessageData{
			Id:        msg.Id,
			ChannelId: msg.ChannelId,
			Content:   msg.Content,
			AuthorId:  msg.AuthorId,
			Reactions: msg.Reactions,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		}
	}
	return exported, nil
}

type ExportableChatMessageData struct {
	Id        entity2.ChannelMessageId          `json:"id"`
	ChannelId channel_entity.ChannelId          `json:"channelId"`
	Content   string                            `json:"content"`
	AuthorId  entity.UserId                     `json:"authorId"`
	Reactions []*entity2.ChannelMessageReaction `json:"reactions,omitempty"`
	CreatedAt time.Time                         `json:"createdAt"`
	UpdatedAt time.Time                         `json:"updatedAt"`
}
