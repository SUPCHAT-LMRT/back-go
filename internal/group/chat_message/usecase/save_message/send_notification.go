package save_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/usecase/send_notification"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/group_info"
	"github.com/supchat-lmrt/back-go/internal/group/usecase/list_members_users"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	uberdig "go.uber.org/dig"
)

type SendNotificationObserverDeps struct {
	uberdig.In
	GetUserByIdUseCase             *get_by_id.GetUserByIdUseCase
	GetGroupInfoUseCase            *group_info.GetGroupInfoUseCase
	ListGroupMembersUseCase        *list_members.ListGroupMembersUseCase
	SendMessageNotificationUseCase *send_notification.SendMessageNotificationUseCase
	Logger                         logger.Logger
}

type SendNotificationObserver struct {
	deps SendNotificationObserverDeps
}

func NewSendNotificationObserver(deps SendNotificationObserverDeps) MessageSavedObserver {
	return &SendNotificationObserver{deps: deps}
}

func (o SendNotificationObserver) NotifyMessageSaved(msg *entity.GroupChatMessage) {
	sender, err := o.deps.GetUserByIdUseCase.Execute(context.Background(), msg.AuthorId)
	if err != nil {
		o.deps.Logger.Error().
			Str("group_message_id", msg.Id.String()).
			Err(err).Msg("échec lors de la récupération de l'expéditeur")
		return
	}

	_, err = o.deps.GetGroupInfoUseCase.Execute(context.Background(), msg.GroupId)
	if err != nil {
		o.deps.Logger.Error().
			Str("group_id", msg.GroupId.String()).
			Err(err).Msg("échec lors de la récupération du groupe")
		return
	}

	// Récupérer tous les membres du groupe
	groupMembers, err := o.deps.ListGroupMembersUseCase.Execute(
		context.Background(),
		msg.GroupId,
	)
	if err != nil {
		o.deps.Logger.Error().
			Str("group_message_id", msg.Id.String()).
			Str("group_id", msg.GroupId.String()).
			Err(err).Msg("échec lors de la récupération des membres du groupe")
		return
	}

	// Envoyer une notification à chaque destinataire sauf l'auteur
	for _, member := range groupMembers {
		if member.UserId == msg.AuthorId {
			continue
		}
		err = o.deps.SendMessageNotificationUseCase.Execute(context.Background(), send_notification.SendMessageNotificationRequest{
			Content:    msg.Content,
			SenderName: sender.FullName(),
			GroupId:    msg.GroupId,
			SenderId:   msg.AuthorId,
			MessageId:  msg.Id,
			ReceiverId: member.UserId,
		})
		if err != nil {
			o.deps.Logger.Error().
				Str("group_message_id", msg.Id.String()).
				Str("recipient_id", member.Id.String()).
				Err(err).Msg("échec lors de l'envoi de la notification")
			continue
		}
	}
}
