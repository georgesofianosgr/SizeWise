package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

var _ Storage = (*S3)(nil)

type S3 struct {
	Base
	AccessKey string     `json:"accessKey"`
	SecretKey string     `json:"secretKey"`
	Bucket    string     `json:"bucket"`
	Region    string     `json:"region"`
	BasePath  string     `json:"basePath"`
	client    *s3.Client `json:"-"`
}

func NewS3(id, basePath, awsKey, awsSecret, bucket, region string) (S3, error) {
	path := basePath
	endsWithSlash := strings.HasSuffix(basePath, "/")
	if path != "" && !endsWithSlash {
		path = basePath + "/"
	}

	customCredentials := aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(awsKey, awsSecret, ""))
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(customCredentials),
	)
	if err != nil {
		return S3{}, err
	}

	client := s3.NewFromConfig(cfg)

	return S3{
		Base: Base{
			ID:   id,
			Type: "s3",
		},
		AccessKey: awsKey,
		SecretKey: awsSecret,
		Bucket:    bucket,
		Region:    region,
		BasePath:  path,
		client:    client,
	}, nil
}

func (s S3) GetID() string {
	return s.ID
}

func (s S3) GetType() string {
	return s.Type
}

func (s S3) ShouldCache() bool {
	return s.Cache
}

func (s S3) EntryExists(path string) (bool, error) {
	key := s.BasePath + path
	headObjectInput := &s3.HeadObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	_, err := s.client.HeadObject(context.TODO(), headObjectInput)
	if err != nil {
		var ErrNotFound *types.NotFound
		if errors.As(err, &ErrNotFound) {
			return false, nil
		}
		return false, fmt.Errorf("unable to check if file exists, %w", err)
	}

	return true, nil
}

func (s S3) WriteEntry(path string, reader io.Reader) error {
	key := s.BasePath + path
	putObjectInput := &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   reader,
	}

	_, err := s.client.PutObject(context.TODO(), putObjectInput)
	if err != nil {
		return fmt.Errorf("unable to upload file, %w", err)
	}

	return nil
}

func (s S3) ReadEntry(path string, writer io.Writer) error {
	key := s.BasePath + path
	getObjectInput := &s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
	}
	result, err := s.client.GetObject(context.TODO(), getObjectInput)
	if err != nil {
		var ErrNotFound *types.NotFound
		if errors.As(err, &ErrNotFound) {
			return ErrEntryNotFound
		}
		return fmt.Errorf("unable to download file, %w", err)
	}
	defer result.Body.Close()
	_, err = io.Copy(writer, result.Body)
	return err
}

func (s *S3) UnmarshalJSON(data []byte) error {
	var tmpData struct {
		Base
		AccessKey string     `json:"accessKey"`
		SecretKey string     `json:"secretKey"`
		Bucket    string     `json:"bucket"`
		Region    string     `json:"region"`
		BasePath  string     `json:"basePath"`
		client    *s3.Client `json:"-"`
	}

	err := json.Unmarshal(data, &tmpData)
	if err != nil {
		return err
	}
	*s, err = NewS3(tmpData.ID, tmpData.BasePath, tmpData.AccessKey, tmpData.SecretKey, tmpData.Bucket, tmpData.Region)
	return err
}
