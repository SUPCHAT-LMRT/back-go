package update_banner

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supchat-lmrt/back-go/internal/s3"
	"github.com/supchat-lmrt/back-go/internal/workspace/entity"
)

type S3UpdateWorkspaceBannerStrategy struct {
	client *s3.S3Client
}

func NewS3UpdateWorkspaceBannerStrategy(client *s3.S3Client) UpdateWorkspaceBannerStrategy {
	return &S3UpdateWorkspaceBannerStrategy{client: client}
}

func (s S3UpdateWorkspaceBannerStrategy) Handle(
	ctx context.Context,
	workspaceId entity.WorkspaceId,
	imageReader io.Reader,
	contentType string,
) error {
	_, err := s.client.Client.PutObject(ctx, &aws_s3.PutObjectInput{
		Key:         aws.String(workspaceId.String()),
		Bucket:      aws.String("workspaces-banners"),
		Body:        imageReader,
		ContentType: aws.String(contentType),
	})
	return err
}
