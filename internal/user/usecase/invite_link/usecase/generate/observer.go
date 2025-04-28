package generate

import "github.com/supchat-lmrt/back-go/internal/user/usecase/invite_link/entity"

type GenerateInviteLinkObserver interface {
	NotifyInviteLinkGenerated(inviteLink *entity.InviteLink, link string)
}
