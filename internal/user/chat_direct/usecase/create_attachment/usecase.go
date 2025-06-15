package create_attachment

import (
	"context"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"github.com/supchat-lmrt/back-go/internal/user/chat_direct/repository"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
	"io"
)

type CreateChatDirectAttachmentUseCaseDeps struct {
	uberdig.In
	ChatDirectRepository repository.ChatDirectRepository
	UploadFileStrategy   UploadFileStrategy
	Observers            []CreateAttachmentObserver `group:"create_chat_direct_attachment_observer"`
}

type CreateChatDirectAttachmentUseCase struct {
	deps CreateChatDirectAttachmentUseCaseDeps
}

func NewCreateChatDirectAttachmentUseCase(deps CreateChatDirectAttachmentUseCaseDeps) *CreateChatDirectAttachmentUseCase {
	return &CreateChatDirectAttachmentUseCase{deps: deps}
}

func (u *CreateChatDirectAttachmentUseCase) Execute(ctx context.Context, input *CreateChatDirectAttachmentInput) (*chat_direct_entity.ChatDirect, error) {
	attachmentId := chat_direct_entity.ChatDirectAttachmentId(bson.NewObjectID().Hex())

	attachment := &chat_direct_entity.ChatDirectAttachment{
		Id:       attachmentId,
		FileName: input.FileName,
	}

	err := u.deps.UploadFileStrategy.Handle(ctx, attachmentId, input.File, input.ContentType)
	if err != nil {
		return nil, err
	}

	chatDirect := chat_direct_entity.ChatDirect{
		SenderId:    input.SenderUserId,
		User1Id:     input.SenderUserId,
		User2Id:     input.OtherUserId,
		Attachments: []*chat_direct_entity.ChatDirectAttachment{attachment},
	}
	err = u.deps.ChatDirectRepository.Create(ctx, &chatDirect)
	if err != nil {
		return nil, err
	}

	// Notify observers about the new attachment
	for _, observer := range u.deps.Observers {
		observer.NotifyAttachmentCreated(&chatDirect)
	}

	return &chatDirect, nil
}

type CreateChatDirectAttachmentInput struct {
	SenderUserId user_entity.UserId
	OtherUserId  user_entity.UserId
	File         io.ReadSeekCloser
	FileName     string
	ContentType  string
}
