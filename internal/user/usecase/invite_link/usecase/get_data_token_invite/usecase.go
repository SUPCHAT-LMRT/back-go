package get_data_token_invite

import (
	"context"

	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/repository"
)

type GetInviteLinkDataUseCase struct {
	repository repository.InviteLinkRepository
}

func NewGetInviteLinkDataUseCase(
	linkRepository repository.InviteLinkRepository,
) *GetInviteLinkDataUseCase {
	return &GetInviteLinkDataUseCase{repository: linkRepository}
}

func (u *GetInviteLinkDataUseCase) GetInviteLinkData(
	ctx context.Context,
	token string,
) (*entity.InviteLink, error) {
	return u.repository.GetInviteLinkData(ctx, token)
}
