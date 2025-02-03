package repository

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
)

type ChatDirectRepository interface {
	Create(ctx context.Context, chatDirect *entity.ChatDirect) error
	ListRecentGroups(ctx context.Context) ([]*entity.ChatDirect, error)
}
