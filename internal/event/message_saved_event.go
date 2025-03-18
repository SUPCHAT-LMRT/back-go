package event

import "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"

const (
	DirectChatMessageSavedEventType EventType = "direct_chat_message_saved"
)

type DirectChatMessageSavedEvent struct {
	Message *entity.ChatDirect
}

func (e DirectChatMessageSavedEvent) Type() EventType {
	return DirectChatMessageSavedEventType
}
