package create_attachment

import (
	"context"
	chat_direct_entity "github.com/supchat-lmrt/back-go/internal/user/chat_direct/entity"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supchat-lmrt/back-go/internal/s3"
)

type S3FileUploadStrategy struct {
	client *s3.S3Client
}

func NewS3FileUploadStrategy(client *s3.S3Client) UploadFileStrategy {
	return &S3FileUploadStrategy{client: client}
}

func (s S3FileUploadStrategy) Handle(
	ctx context.Context,
	attachmentId chat_direct_entity.ChatDirectAttachmentId,
	fileReader io.Reader,
	contentType string,
) error {
	_, err := s.client.Client.PutObject(ctx, &aws_s3.PutObjectInput{
		Key:         aws.String(attachmentId.String()),
		Bucket:      aws.String("chat-direct-attachments"),
		Body:        fileReader,
		ContentType: aws.String(contentType),
	})
	return err
}
