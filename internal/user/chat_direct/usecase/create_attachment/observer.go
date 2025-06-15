package create_attachment

import (
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
)

type CreateAttachmentObserver interface {
	NotifyAttachmentCreated(message *chat_direct_entity.ChatDirect)
}
