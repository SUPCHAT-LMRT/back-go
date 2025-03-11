package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type ChatDirectRepository interface {
	Create(ctx context.Context, chatDirect *entity.ChatDirect) error
	ListRecentChats(ctx context.Context) ([]*entity.ChatDirect, error)
	// ListByUser returns all direct chats between user1 and user2
	ListByUser(ctx context.Context, user1Id, user2Id user_entity.UserId) ([]*entity.ChatDirect, error)
}
