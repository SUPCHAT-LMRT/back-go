package save_message

import (
	"context"
	"time"

	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type SaveMessageInput struct {
	Id        entity.GroupChatMessageId
	GroupId   group_entity.GroupId
	SenderId  user_entity.UserId
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SaveGroupChatMessageUseCase struct {
	repository repository.ChatMessageRepository
}

func NewSaveGroupChatMessageUseCase(
	repository repository.ChatMessageRepository,
) *SaveGroupChatMessageUseCase {
	return &SaveGroupChatMessageUseCase{repository: repository}
}

func (u *SaveGroupChatMessageUseCase) Execute(
	ctx context.Context,
	msg *entity.GroupChatMessage,
) error {
	if err := u.repository.Create(ctx, msg); err != nil {
		return err
	}

	// TODO impl meilisearch

	return nil
}
