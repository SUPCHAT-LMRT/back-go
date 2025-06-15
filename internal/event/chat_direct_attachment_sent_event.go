package event

import (
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
)

const (
	ChatDirectAttachmentSentEventType EventType = "chat_direct_attachment_sent"
)

type ChatDirectAttachmentSentEvent struct {
	ChatDirect *chat_direct_entity.ChatDirect
}

func (e ChatDirectAttachmentSentEvent) Type() EventType {
	return ChatDirectAttachmentSentEventType
}
