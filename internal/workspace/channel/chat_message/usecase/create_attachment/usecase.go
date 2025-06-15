package create_attachment

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/repository"
	channel_entity "github.com/supchat-lmrt/back-go/internal/workspace/channel/entity"
	workspace_entity "github.com/supchat-lmrt/back-go/internal/workspace/entity"
	workspace_member_entity "github.com/supchat-lmrt/back-go/internal/workspace/member/entity"
	"go.mongodb.org/mongo-driver/v2/bson"
	uberdig "go.uber.org/dig"
	"io"
)

type CreateChannelMessageAttachmentUseCaseDeps struct {
	uberdig.In
	ChannelMessageRepository repository.ChannelMessageRepository
	UploadFileStrategy       UploadFileStrategy
	Observers                []CreateAttachmentObserver `group:"create_channel_messages_attachment_observer"`
}

type CreateChannelMessageAttachmentUseCase struct {
	deps CreateChannelMessageAttachmentUseCaseDeps
}

func NewCreateChatDirectAttachmentUseCase(deps CreateChannelMessageAttachmentUseCaseDeps) *CreateChannelMessageAttachmentUseCase {
	return &CreateChannelMessageAttachmentUseCase{deps: deps}
}

func (u *CreateChannelMessageAttachmentUseCase) Execute(ctx context.Context, input *CreateChannelMessageAttachmentInput) (*entity.ChannelMessage, error) {
	attachmentId := entity.ChannelMessageAttachmentId(bson.NewObjectID().Hex())

	attachment := &entity.ChannelMessageAttachment{
		Id:       attachmentId,
		FileName: input.FileName,
	}

	err := u.deps.UploadFileStrategy.Handle(ctx, attachmentId, input.File, input.ContentType)
	if err != nil {
		return nil, err
	}

	channelMessage := entity.ChannelMessage{
		ChannelId:   input.ChannelId,
		AuthorId:    input.SenderWorkspaceMember.UserId,
		Attachments: []*entity.ChannelMessageAttachment{attachment},
	}
	err = u.deps.ChannelMessageRepository.Create(ctx, &channelMessage)
	if err != nil {
		return nil, err
	}

	// Notify observers about the new attachment
	for _, observer := range u.deps.Observers {
		observer.NotifyAttachmentCreated(input.WorkspaceId, input.SenderWorkspaceMember.Id, &channelMessage)
	}

	return &channelMessage, nil
}

type CreateChannelMessageAttachmentInput struct {
	WorkspaceId           workspace_entity.WorkspaceId
	ChannelId             channel_entity.ChannelId
	SenderWorkspaceMember *workspace_member_entity.WorkspaceMember
	File                  io.ReadSeekCloser
	FileName              string
	ContentType           string
}
