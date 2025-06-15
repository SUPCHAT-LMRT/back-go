package save_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/search/message"
	uberdig "go.uber.org/dig"
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

type SaveGroupChatMessageUseCaseDeps struct {
	uberdig.In
	ChatMessageRepository    repository.ChatMessageRepository
	SearchMessageSyncManager message.SearchMessageSyncManager
	Observers                []MessageSavedObserver `group:"send_groupmessage_notification_channel"`
}

type SaveGroupChatMessageUseCase struct {
	deps SaveGroupChatMessageUseCaseDeps
}

func NewSaveGroupChatMessageUseCase(
	deps SaveGroupChatMessageUseCaseDeps,
) *SaveGroupChatMessageUseCase {
	return &SaveGroupChatMessageUseCase{deps: deps}
}

func (u *SaveGroupChatMessageUseCase) Execute(
	ctx context.Context,
	msg *entity.GroupChatMessage,
) error {
	if err := u.deps.ChatMessageRepository.Create(ctx, msg); err != nil {
		return err
	}

	err := u.deps.SearchMessageSyncManager.AddMessage(ctx, &message.SearchMessage{
		Id:      msg.Id.String(),
		Content: msg.Content,
		Kind:    message.SearchMessageGroupMessage,
		Data: message.SearchMessageGroupData{
			GroupId: msg.GroupId,
		},
		AuthorId:  msg.AuthorId,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
	})
	if err != nil {
		return err
	}
	// Notifier les observateurs
	for _, observer := range u.deps.Observers {
		observer.NotifyMessageSaved(msg)
	}

	return nil
}
