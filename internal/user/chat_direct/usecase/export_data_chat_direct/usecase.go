package export_data_chat_direct

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ExportDataChatDirectUseCase struct {
	repo repository.ChatDirectRepository
}

func NewExportDataChatDirectUseCase(repo repository.ChatDirectRepository) *ExportDataChatDirectUseCase {
	return &ExportDataChatDirectUseCase{repo: repo}
}

func (uc *ExportDataChatDirectUseCase) Execute(ctx context.Context, userId user_entity.UserId) ([]ExportableChatDirectData, error) {
	messages, err := uc.repo.ListAllMessagesByUser(ctx, userId)
	if err != nil {
		return nil, err
	}

	exported := make([]ExportableChatDirectData, len(messages))
	for i, msg := range messages {
		exported[i] = ExportableChatDirectData{
			Id:        msg.Id.String(),
			SenderId:  msg.SenderId.String(),
			User1Id:   msg.User1Id.String(),
			User2Id:   msg.User2Id.String(),
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
		}
	}
	return exported, nil
}

type ExportableChatDirectData struct {
	Id        string    `json:"id"`
	SenderId  string    `json:"senderId"`
	User1Id   string    `json:"user1Id"`
	User2Id   string    `json:"user2Id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
