package create_attachment

import (
	"context"
	"github.com/supchat-lmrt/back-go/internal/workspace/channel/chat_message/entity"
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
	attachmentId entity.ChannelMessageAttachmentId,
	fileReader io.Reader,
	contentType string,
) error {
	_, err := s.client.Client.PutObject(ctx, &aws_s3.PutObjectInput{
		Key:         aws.String(attachmentId.String()),
		Bucket:      aws.String("channels-attachments"),
		Body:        fileReader,
		ContentType: aws.String(contentType),
	})
	return err
}
