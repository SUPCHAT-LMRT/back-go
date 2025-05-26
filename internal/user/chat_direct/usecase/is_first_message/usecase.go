package is_first_message

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
)

type IsFirstMessageUseCase struct {
	repository repository.ChatDirectRepository
}

func NewIsFirstMessageUseCase(repository repository.ChatDirectRepository) *IsFirstMessageUseCase {
	return &IsFirstMessageUseCase{repository: repository}
}

func (u *IsFirstMessageUseCase) Execute(
	ctx context.Context,
	user1Id, user2Id user_entity.UserId,
) (bool, error) {
	return u.repository.IsFirstMessage(ctx, user1Id, user2Id)
}
