package save_message

import "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"

type MessageSavedObserver interface {
	NotifyMessageSaved(msg *entity.ChatDirect)
}
