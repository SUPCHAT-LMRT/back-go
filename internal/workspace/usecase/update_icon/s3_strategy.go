package update_icon

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supchat-lmrt/back-go/internal/s3"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
	"io"
)

type S3UpdateWorkspaceIconStrategy struct {
	client *s3.S3Client
}

func NewS3UpdateWorkspaceIconStrategy(client *s3.S3Client) UpdateWorkspaceIconStrategy {
	return &S3UpdateWorkspaceIconStrategy{client: client}
}

func (s S3UpdateWorkspaceIconStrategy) Handle(ctx context.Context, workspaceId entity.WorkspaceId, imageReader io.Reader, contentType string) error {
	_, err := s.client.Client.PutObject(ctx, &aws_s3.PutObjectInput{
		Key:         aws.String(workspaceId.String()),
		Bucket:      aws.String("workspaces-icons"),
		Body:        imageReader,
		ContentType: aws.String(contentType),
	})
	return err
}
