package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"time"
)

type ChatDirectRepository interface {
	Create(ctx context.Context, chatDirect *entity.ChatDirect) error
	ListRecentChats(ctx context.Context, userId user_entity.UserId) ([]*entity.ChatDirect, error)
	IsFirstMessage(ctx context.Context, user1Id, user2Id user_entity.UserId) (bool, error)
	// ListByUser returns all direct chats between user1 and user2
	ListByUser(ctx context.Context, user1Id, user2Id user_entity.UserId, params ListByUserQueryParams) ([]*entity.ChatDirect, error)
	ToggleReaction(ctx context.Context, messageId entity.ChatDirectId, userId user_entity.UserId, reaction string) (added bool, err error)
}

type ListByUserQueryParams struct {
	Limit           int
	Before          time.Time
	After           time.Time
	AroundMessageId entity.ChatDirectId
}
