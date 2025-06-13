package websocket

import (
	group_chat_entity "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"github.com/supchat-lmrt/back-go/internal/websocket/messages/inbound"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
)

type SendChannelMessageObserver interface {
	OnSendMessage(
		message *inbound.InboundSendMessageToChannel,
		messageId entity.ChannelMessageId,
		userId user_entity.UserId,
	)
}

type SendDirectMessageObserver interface {
	OnSendMessage(
		message *inbound.InboundSendDirectMessage,
		messageId chat_direct_entity.ChatDirectId,
		userId user_entity.UserId,
	)
}

type SendGroupMessageObserver interface {
	OnSendMessage(
		message *inbound.InboundSendGroupMessage,
		messageId group_chat_entity.GroupChatMessageId,
		userId user_entity.UserId,
	)
}
