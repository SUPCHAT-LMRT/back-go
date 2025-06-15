package save_message

import (
	entity2 "github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
)

type MessageSavedObserver interface {
	NotifyMessageSaved(msg *entity2.GroupChatMessage)
}
