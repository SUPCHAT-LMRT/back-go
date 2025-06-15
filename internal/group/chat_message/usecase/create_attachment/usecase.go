package create_attachment

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/group/chat_message/entity"
	group_chat_message_repository "github.com/supchat-lmrt/back-go/internal/group/chat_message/repository"
	group_entity "github.com/supchat-lmrt/back-go/internal/group/entity"
	user_entity "github.com/supchat-lmrt/back-go/internal/user/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
	"io"
)

type CreateGroupMessageAttachmentUseCaseDeps struct {
	uberdig.In
	ChatMessageRepository group_chat_message_repository.ChatMessageRepository
	UploadFileStrategy    UploadFileStrategy
	Observers             []CreateAttachmentObserver `group:"create_group_attachment_observer"`
}

type CreateGroupAttachmentUseCase struct {
	deps CreateGroupMessageAttachmentUseCaseDeps
}

func NewCreateGroupAttachmentUseCase(deps CreateGroupMessageAttachmentUseCaseDeps) *CreateGroupAttachmentUseCase {
	return &CreateGroupAttachmentUseCase{deps: deps}
}

func (u *CreateGroupAttachmentUseCase) Execute(ctx context.Context, input *CreateGroupAttachmentInput) (*entity.GroupChatMessage, error) {
	attachmentId := entity.GroupChatAttachmentId(bson.NewObjectID().Hex())

	attachment := &entity.GroupChatMessageAttachment{
		Id:       attachmentId,
		FileName: input.FileName,
	}

	err := u.deps.UploadFileStrategy.Handle(ctx, attachmentId, input.File, input.ContentType)
	if err != nil {
		return nil, err
	}

	groupChatMessage := entity.GroupChatMessage{
		GroupId:     input.GroupId,
		AuthorId:    input.SenderUserId,
		Attachments: []*entity.GroupChatMessageAttachment{attachment},
	}
	err = u.deps.ChatMessageRepository.Create(ctx, &groupChatMessage)
	if err != nil {
		return nil, err
	}

	// Notify observers about the new attachment
	for _, observer := range u.deps.Observers {
		observer.NotifyAttachmentCreated(&groupChatMessage)
	}

	return &groupChatMessage, nil
}

type CreateGroupAttachmentInput struct {
	GroupId      group_entity.GroupId
	SenderUserId user_entity.UserId
	File         io.ReadSeekCloser
	FileName     string
	ContentType  string
}
