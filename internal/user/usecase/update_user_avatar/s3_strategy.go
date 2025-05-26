package update_user_avatar

import (
	"context"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/supchat-lmrt/back-go/internal/s3"
	"github.com/supchat-lmrt/back-go/internal/user/entity"
)

type S3UpdateUserAvatarStrategy struct {
	client *s3.S3Client
}

func NewS3UpdateUserAvatarStrategy(client *s3.S3Client) UpdateUserAvatarStrategy {
	return &S3UpdateUserAvatarStrategy{client: client}
}

func (s S3UpdateUserAvatarStrategy) Handle(
	ctx context.Context,
	userId entity.UserId,
	imageReader io.Reader,
	contentType string,
) error {
	_, err := s.client.Client.PutObject(ctx, &aws_s3.PutObjectInput{
		Key:         aws.String(userId.String()),
		Bucket:      aws.String("users-avatars"),
		Body:        imageReader,
		ContentType: aws.String(contentType),
	})
	return err
}
