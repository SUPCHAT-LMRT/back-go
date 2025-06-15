package event

import "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"

const (
	GroupAttachmentSentEventType EventType = "group_attachment_sent"
)

type GroupAttachmentSentEvent struct {
	GroupChatMessage *entity.GroupChatMessage
}

func (e GroupAttachmentSentEvent) Type() EventType {
	return GroupAttachmentSentEventType
}
