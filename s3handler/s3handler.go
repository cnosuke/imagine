package s3handler

import (
	"context"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/cnosuke/imagine/entity"
	"github.com/satori/go.uuid"
)

type S3Handler struct {
	svc                 *s3.S3
	ctx                 context.Context
	bucket              string
	keyPrefix           string
	defaultPresignedTTL time.Duration
}

func NewS3Handler(ctx context.Context, awsRegion string, bucket string, keyPrefix string, defaultPresignedTTL int) *S3Handler {
	sess := session.Must(
		session.NewSession(&aws.Config{Region: aws.String(awsRegion)}),
	)

	return &S3Handler{
		svc:                 s3.New(sess),
		ctx:                 ctx,
		bucket:              bucket,
		keyPrefix:           keyPrefix,
		defaultPresignedTTL: time.Duration(defaultPresignedTTL) * time.Second,
	}
}

func (s *S3Handler) CreatePresignedPostUrl(filename string, contentType string) (*entity.PresignedPostUrl, error) {
	return s.CreatePresignedPostUrlWithTTL(filename, contentType, s.defaultPresignedTTL)
}

func (s *S3Handler) CreatePresignedPostUrlWithTTL(filename string, contentType string, ttl time.Duration) (*entity.PresignedPostUrl, error) {
	id := uuid.NewV4().String()

	key := filepath.Join(s.keyPrefix, id, filename)

	req, _ := s.svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	})

	str, err := req.Presign(ttl)
	if err != nil {
		return &entity.PresignedPostUrl{}, err
	}

	return &entity.PresignedPostUrl{
		Url:         str,
		Ttl:         ttl,
		Key:         key,
		Id:          id,
		Filename:    filename,
		ContentType: contentType,
	}, nil
}
