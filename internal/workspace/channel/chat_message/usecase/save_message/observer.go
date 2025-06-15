package save_message

import "github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"

type MessageSavedObserver interface {
	NotifyMessageSaved(msg *entity.ChannelMessage)
}
