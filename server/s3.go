package server

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/tomnlittle/gocv-server/cache"
)

// AwsConfig configuration
type AwsConfig struct {
	S3Downloader *s3manager.Downloader
	Cache        *cache.ImageCache
}

// NewAwsConfig returns a new AWS configuration
func NewAwsConfig(mc *cache.ImageCache) (*AwsConfig, error) {

	// intialise aws session
	sess, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})

	if err != nil {
		return nil, err
	}

	return &AwsConfig{
		S3Downloader: s3manager.NewDownloader(sess),
		Cache:        mc,
	}, nil
}

// GetObject returns an s3 object
func (a *AwsConfig) GetObject(bucket, key string) ([]byte, error) {

	// check if the image is in the cache
	hash := a.Cache.GenerateHash(bucket, key)
	bytes, err := a.Cache.GetBytes(hash)

	if err != nil {
		return nil, err
	}

	if bytes != nil {
		return bytes, nil
	}

	params := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	buf := aws.NewWriteAtBuffer([]byte{})
	if _, err := a.S3Downloader.Download(buf, params); err != nil {
		return nil, err
	}

	if err = a.Cache.AddBytes(hash, buf.Bytes()); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
