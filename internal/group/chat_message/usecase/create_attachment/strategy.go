package create_attachment

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	"io"
)

type UploadFileStrategy interface {
	Handle(
		ctx context.Context,
		attachmentId entity.GroupChatAttachmentId,
		fileReader io.Reader,
		contentType string,
	) error
}
