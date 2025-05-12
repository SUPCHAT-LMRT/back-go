package get_list_invite_link

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/repository"
)

type GetListInviteLinkUseCase struct {
	repository repository.InviteLinkRepository
}

func NewGetListInviteLinkUseCase(repo repository.InviteLinkRepository) *GetListInviteLinkUseCase {
	return &GetListInviteLinkUseCase{repository: repo}
}

func (u *GetListInviteLinkUseCase) Execute(ctx context.Context) ([]*entity.InviteLink, error) {
	return u.repository.GetAllInviteLinks(ctx)
}
