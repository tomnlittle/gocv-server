package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// AwsConfig configuration
type AwsConfig struct {
	S3Downloader *s3manager.Downloader
}

// NewAwsConfig returns a new AWS configuration
func NewAwsConfig() (*AwsConfig, error) {

	// intialise aws session
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, err
	}

	return &AwsConfig{
		S3Downloader: s3manager.NewDownloader(sess),
	}, nil
}

// GetObject returns an s3 object
func (a *AwsConfig) GetObject(bucket, key string) ([]byte, error) {

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	buf := aws.NewWriteAtBuffer([]byte{})
	if _, err := a.S3Downloader.Download(buf, params); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
