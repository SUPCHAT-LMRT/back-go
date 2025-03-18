package get_data_token_invite

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/invite_link_workspace/repository"
)

type GetInviteLinkDataUseCase struct {
	repository repository.InviteLinkRepository
}

func NewGetInviteLinkDataUseCase(linkRepository repository.InviteLinkRepository) *GetInviteLinkDataUseCase {
	return &GetInviteLinkDataUseCase{repository: linkRepository}
}

func (u *GetInviteLinkDataUseCase) GetInviteLinkData(ctx context.Context, token string) (*entity.InviteLink, error) {
	inviteLinkData, err := u.repository.GetInviteLinkData(ctx, token)
	if err != nil {
		return nil, err
	}

	return inviteLinkData, nil
}
