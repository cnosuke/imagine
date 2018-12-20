package s3handler

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/satori/go.uuid"
	"path/filepath"
	"time"
)

type S3Handler struct {
	svc *s3.S3
	ctx context.Context
	bucket string
	keyPrefix string
}

type PresignedPostUrl struct {
	url string
	ttl time.Duration
	key string
	id string
	filename string
	contentType string
}

var AwsRegion = "ap-northeast-1"
var DefaultTTL = 5 * time.Minute

func NewS3Handler(ctx context.Context, bucket string, keyPrefix string) *S3Handler {
	sess := session.Must(
		session.NewSession(&aws.Config{Region: aws.String(AwsRegion)}),
	)

	return &S3Handler{
		svc: s3.New(sess),
		ctx: ctx,
		bucket: bucket,
		keyPrefix: keyPrefix,
	}
}

func (s *S3Handler) CreatePresignedPostUrl(filename string, contentType string) (*PresignedPostUrl, error) {
	return s.CreatePresignedPostUrlWithTTL(filename, contentType, DefaultTTL)
}


func (s *S3Handler) CreatePresignedPostUrlWithTTL(filename string, contentType string, ttl time.Duration) (*PresignedPostUrl, error) {
	id := uuid.Must(uuid.NewV4()).String()

	key := filepath.Join(s.keyPrefix, id, filename)

	req, _ := s.svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
		ContentType: aws.String(contentType),
	})

	str, err := req.Presign(ttl)
	if err != nil {
		return &PresignedPostUrl{}, err
	}

	return &PresignedPostUrl{
		url: str,
		ttl: ttl,
		key: key,
		id: id,
		filename: filename,
		contentType: contentType,
	}, nil
}


