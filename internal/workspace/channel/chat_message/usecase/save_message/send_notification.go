package save_message

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/logger"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/send_notification"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"
	"github.com/supchat-lmrt/back-go/internal/workspace/member/usecase/list_workspace_members"
	uberdig "go.uber.org/dig"
)

type SendNotificationObserverDeps struct {
	uberdig.In
	GetUserByIdUseCase             *get_by_id.GetUserByIdUseCase
	GetChannelUseCase              *get_channel.GetChannelUseCase
	ListWorkspaceMembersUseCase    *list_workspace_members.ListWorkspaceMembersUseCase
	SendMessageNotificationUseCase *send_notification.SendMessageNotificationUseCase
	Logger                         logger.Logger
}

type SendNotificationObserver struct {
	deps SendNotificationObserverDeps
}

func NewSendNotificationObserver(deps SendNotificationObserverDeps) MessageSavedObserver {
	return &SendNotificationObserver{deps: deps}
}

func (o SendNotificationObserver) NotifyMessageSaved(msg *entity.ChannelMessage) {
	sender, err := o.deps.GetUserByIdUseCase.Execute(context.Background(), msg.AuthorId)
	if err != nil {
		o.deps.Logger.Error().
			Str("channel_message_id", msg.Id.String()).
			Err(err).Msg("failed to get sender user")
		return
	}

	channel, err := o.deps.GetChannelUseCase.Execute(context.Background(), msg.ChannelId)
	if err != nil {
		o.deps.Logger.Error().
			Str("channel_message_id", msg.ChannelId.String()).
			Err(err).Msg("failed to get channel")
		return
	}
	var recipients []string
	if channel.IsPrivate {
		recipients = channel.Members
	} else {
		// Pour un canal public, récupérer tous les membres du workspace
		// On met une limite élevée pour récupérer tous les membres en une fois
		_, workspaceMembers, err := o.deps.ListWorkspaceMembersUseCase.Execute(
			context.Background(),
			channel.WorkspaceId,
			1000, // limite maximum par page
			1,    // première page
		)
		if err != nil {
			o.deps.Logger.Error().
				Str("channel_message_id", msg.Id.String()).
				Str("workspace_id", channel.WorkspaceId.String()).
				Err(err).Msg("failed to get workspace members")
			return
		}
		for _, member := range workspaceMembers {
			recipients = append(recipients, string(member.UserId))
		}
	}

	// Envoyer une notification à chaque destinataire sauf l'auteur
	for _, recipientId := range recipients {
		if recipientId == string(msg.AuthorId) {
			continue
		}
		err = o.deps.SendMessageNotificationUseCase.Execute(context.Background(), send_notification.SendMessageNotificationRequest{
			Content:     msg.Content,
			SenderName:  sender.FullName(),
			SenderId:    msg.AuthorId,
			MessageId:   msg.Id,
			ChannelId:   msg.ChannelId,
			WorkspaceId: channel.WorkspaceId,
			ReceiverId:  user_entity.UserId(recipientId),
		})
		if err != nil {
			o.deps.Logger.Error().
				Str("channel_message_id", msg.Id.String()).
				Str("recipient_id", recipientId).
				Err(err).Msg("failed to send message notification")
			continue
		}
	}
}
