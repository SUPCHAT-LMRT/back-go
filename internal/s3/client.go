package s3

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	transport "github.com/aws/smithy-go/endpoints"
)

type S3Client struct {
	Client *s3.Client
}

func NewS3Client() (*S3Client, error) {
	parsedUrl, err := url.Parse(os.Getenv("S3_ENDPOINT"))
	if err != nil {
		return nil, err
	}

	var endpointResolver s3.EndpointResolverV2
	if os.Getenv("S3_ENDPOINT_MODE") == "path" {
		endpointResolver = &pathResolver{}
	} else if os.Getenv("S3_ENDPOINT_MODE") == "host" {
		endpointResolver = s3.NewDefaultEndpointResolverV2()
	}

	s3Client := s3.New(s3.Options{
		EndpointResolverV2: endpointResolver,
		BaseEndpoint:       aws.String(parsedUrl.String()),
		Credentials: credentials.NewStaticCredentialsProvider(
			os.Getenv("S3_ACCESS_KEY"),
			os.Getenv("S3_SECRET_KEY"),
			"",
		),
		Region: "eu-west-1",
	})

	return &S3Client{Client: s3Client}, nil
}

func (s *S3Client) CreateBucketIfNotExist(
	ctx context.Context,
	bucketName string,
) (bucketsCreated bool, err error) {
	_, _ = s.Client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: &bucketName,
	})

	// Make the bucket objects publicly accessible
	_, _ = s.Client.PutBucketPolicy(ctx, &s3.PutBucketPolicyInput{
		Bucket: &bucketName,
		Policy: aws.String(fmt.Sprintf(`{
    "Version": "2012-10-17",
    "Statement": [
        {
            "Effect": "Allow",
            "Principal": {
                "AWS": [
                    "*"
                ]
            },
            "Action": [
                "s3:GetObject"
            ],
            "Resource": [
                "arn:aws:s3:::%s/*"
            ]
        }
    ]
}`, bucketName)),
	})
	if err != nil {
		var bne *types.BucketAlreadyOwnedByYou
		if errors.As(err, &bne) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// pathResolver is a custom endpoint resolver for the S3 client. It appends the bucket name to the endpoint URL. (for minio)
type pathResolver struct{}

func (r *pathResolver) ResolveEndpoint(
	_ context.Context,
	params s3.EndpointParameters,
) (transport.Endpoint, error) {
	parsedEndpoint, err := url.Parse(*params.Endpoint)
	if err != nil {
		return transport.Endpoint{}, fmt.Errorf("failed to parse endpoint URL: %w", err)
	}
	u := *parsedEndpoint
	u.Path += "/" + *params.Bucket
	return transport.Endpoint{URI: u}, nil
}
