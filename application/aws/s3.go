package aws

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"io"
	"json2xls/config"
	"time"
)

const (
	NotFoundError     = "%s Not Found!"
	FileUploadError   = "AWS s3 file upload error"
	PresignedUrlError = "AWS s3 presigned url request error"
)

type S3FileUploader struct {
	URL        string
	BucketName string
	Region     string
	AccessKey  string
	SecretKey  string
	DisableSSL bool
}

func New() S3FileUploader {
	c := config.New()

	return S3FileUploader{
		URL:        c.AwsURL,
		BucketName: c.AwsBucketName,
		Region:     c.AwsRegion,
		AccessKey:  c.AwsAccessKey,
		SecretKey:  c.AwsSecretKey,
		DisableSSL: c.AwsDisableSSL,
	}
}

func (s3u S3FileUploader) Upload(data io.Reader, key string, signedTTL int64) (uri string, err error) {
	sess, err := s3u.openSession()
	if err != nil {
		return "", err
	}

	svc := s3manager.NewUploader(sess)
	upParams := &s3manager.UploadInput{
		Bucket: aws.String(s3u.BucketName),
		Key:    aws.String(key),
		Body:   data,
	}

	_, err = svc.Upload(upParams)
	if err != nil {
		return "", errors.New(FileUploadError)
	}

	uri, err = s3u.SignedRetrievalURL(key, key, signedTTL)

	if err != nil {
		return "", errors.New(FileUploadError)
	}

	return uri, err

}

func (s3u S3FileUploader) Delete(key string) error {
	sess, err := s3u.openSession()
	if err != nil {
		return err
	}

	s3Client := s3.New(sess)

	req, _ := s3Client.DeleteObjectRequest(&s3.DeleteObjectInput{
		Bucket: aws.String(s3u.BucketName),
		Key:    aws.String(key),
	})
	err = req.Send()
	if err != nil {
		return errors.New(FileUploadError)
	}

	return nil
}

func (s3u S3FileUploader) Ping() error {
	sess, err := s3u.openSession()
	if err != nil {
		return err
	}

	s3Client := s3.New(sess)
	_, err = s3Client.HeadBucket(&s3.HeadBucketInput{Bucket: &s3u.BucketName})
	if err != nil {
		return errors.New(fmt.Sprintf(NotFoundError, "Bucket"))
	}

	return nil
}

func (s3u S3FileUploader) openSession() (*session.Session, error) {
	conf := aws.Config{
		Region:           aws.String(s3u.Region),
		Endpoint:         aws.String(s3u.URL),
		S3ForcePathStyle: aws.Bool(true),
		DisableSSL:       aws.Bool(s3u.DisableSSL),
	}

	if s3u.AccessKey != "" && s3u.SecretKey != "" {
		crds := credentials.Value{AccessKeyID: s3u.AccessKey, SecretAccessKey: s3u.SecretKey}
		creds := credentials.NewStaticCredentialsFromCreds(crds)
		_, err := creds.Get()
		if err != nil {
			return nil, errors.New(FileUploadError)
		}
		conf.Credentials = creds
	}

	sess, err := session.NewSession(&conf)
	if err != nil {
		return nil, errors.New(FileUploadError)
	}

	return sess, nil
}

func (s3u S3FileUploader) SignedRetrievalURL(key string, originalFileName string, signedTTL int64) (uri string, err error) {
	sess, err := s3u.openSession()
	if err != nil {
		return "", err
	}

	s3Client := s3.New(sess)
	req, _ := s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket:                     aws.String(s3u.BucketName),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(fmt.Sprintf("attachment; filename =\"%s\"", originalFileName)),
	})

	uri, err = req.Presign(time.Duration(signedTTL) * time.Minute)
	if err != nil {
		return "", errors.New(PresignedUrlError)
	}

	return uri, nil
}
