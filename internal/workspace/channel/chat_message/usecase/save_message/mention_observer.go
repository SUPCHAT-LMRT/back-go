package save_message

import (
	"context"
	"fmt"
	"github.com/supchat-lmrt/back-go/internal/logger"
	"github.com/supchat-lmrt/back-go/internal/mention/usecase/extract_mentions"
	"github.com/supchat-lmrt/back-go/internal/user/usecase/get_by_id"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/usecase/send_notification"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/usecase/get_channel"

	uberdig "go.uber.org/dig"
)

type GetMentionObserverDeps struct {
	uberdig.In
	ExtractMentionsUseCase         *extract_mentions.ExtractMentionsUseCase
	GetUserByIdUseCase             *get_by_id.GetUserByIdUseCase
	GetChannelUseCase              *get_channel.GetChannelUseCase
	SendMessageNotificationUseCase *send_notification.SendMessageNotificationUseCase
	Logger                         logger.Logger
}

type GetMentionObserver struct {
	deps GetMentionObserverDeps
}

func NewGetMentionObserver(deps GetMentionObserverDeps) MessageSavedObserver {
	return &GetMentionObserver{deps: deps}
}

func (o GetMentionObserver) NotifyMessageSaved(msg *entity.ChannelMessage) {
	fmt.Println("1")
	userIds := o.deps.ExtractMentionsUseCase.Execute(msg.Content)
	if len(userIds) == 0 {
		return
	}
	fmt.Println("2")

	sender, err := o.deps.GetUserByIdUseCase.Execute(context.Background(), msg.AuthorId)
	if err != nil {
		o.deps.Logger.Error().
			Str("channel_message_id", msg.Id.String()).
			Err(err).Msg("failed to get sender user")
		return
	}

	fmt.Println("3")
	channel, err := o.deps.GetChannelUseCase.Execute(context.Background(), msg.ChannelId)
	if err != nil {
		o.deps.Logger.Error().
			Str("channel_message_id", msg.ChannelId.String()).
			Err(err).Msg("failed to get channel")
		return
	}

	fmt.Println("4")
	for _, userId := range userIds {
		mentionnedUserId := userId
		if userId == msg.AuthorId {
			continue
		}

		fmt.Println("5")
		err = o.deps.SendMessageNotificationUseCase.Execute(
			context.Background(),
			send_notification.SendMessageNotificationRequest{
				Content:     msg.Content,
				SenderName:  sender.FullName(),
				SenderId:    msg.AuthorId,
				MessageId:   msg.Id,
				ChannelId:   msg.ChannelId,
				WorkspaceId: channel.WorkspaceId,
				ReceiverId:  mentionnedUserId,
			},
		)
		if err != nil {
			o.deps.Logger.Error().
				Str("channel_message_id", msg.Id.String()).
				Str("mentioned_user_id", string(mentionnedUserId)).
				Err(err).Msg("failed to send mention notification")
			continue
		}
	}
}
