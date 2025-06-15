package create_attachment

import (
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
)

type CreateAttachmentObserver interface {
	NotifyAttachmentCreated(message *entity.GroupChatMessage)
}
