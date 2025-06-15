package create_attachment

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"io"
)

type UploadFileStrategy interface {
	Handle(
		ctx context.Context,
		attachmentId entity.ChannelMessageAttachmentId,
		fileReader io.Reader,
		contentType string,
	) error
}
