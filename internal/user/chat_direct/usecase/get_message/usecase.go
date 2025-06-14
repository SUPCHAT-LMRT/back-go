package get_message

import (
	"context"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
)

type GetDirectChatMessageUseCase struct {
	chatDirectRepository repository.ChatDirectRepository
}

func NewGetDirectChatMessageUseCase(
	chatDirectRepository repository.ChatDirectRepository,
) *GetDirectChatMessageUseCase {
	return &GetDirectChatMessageUseCase{chatDirectRepository: chatDirectRepository}
}

func (u *GetDirectChatMessageUseCase) Execute(
	ctx context.Context,
	msgId chat_direct_entity.ChatDirectId,
) (*chat_direct_entity.ChatDirect, error) {
	return u.chatDirectRepository.GetById(ctx, msgId)
}
