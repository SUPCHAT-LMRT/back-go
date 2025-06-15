package create_attachment

import (
	"context"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"io"
)

type UploadFileStrategy interface {
	Handle(
		ctx context.Context,
		attachmentId chat_direct_entity.ChatDirectAttachmentId,
		fileReader io.Reader,
		contentType string,
	) error
}
