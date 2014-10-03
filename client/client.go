package client

import (
	"errors"
	"fmt"
	"github.com/hailocab/goamz/aws"
	"github.com/hailocab/goamz/s3"
	"regexp"
	"time"
)

type S3Client struct {
	S3Bucket     *s3.Bucket
	AWSRegion    aws.Region
	AWSAuth      aws.Auth
	S3Connection *s3.S3
	Err          error
	BucketName   string
}

func NewClient(bucketName, accessKey, secretKey, awsRegionName string) (*S3Client, error) {

	s3client := &S3Client{}
	s3client.BucketName = bucketName

	// Validate the region exists and fail
	awsRegion, regionExists := aws.Regions[awsRegionName]

	if !regionExists {
		return nil, errors.New(fmt.Sprintf("Region doesn't exist (%s)!\n", awsRegionName))
	} else {
		s3client.AWSRegion = awsRegion
		fmt.Printf("Region: '%s'\n", awsRegionName)
	}

	//fmt.Printf("%s:%s", accessKey, secretKey)
	// Authenticate
	auth, err := aws.GetAuth(accessKey, secretKey, "", time.Time{})

	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to authenticate to AWS (%s)", err))
	}

	s3client.AWSAuth = auth
	s3client.S3Connection = s3.New(s3client.AWSAuth, s3client.AWSRegion)
	s3client.S3Bucket = s3client.S3Connection.Bucket(s3client.BucketName)

	// Checking connection by quering non existant bucket
	if err = s3client.testConnection(); err != nil {
		return nil, err
	}

	return s3client, nil
}

func (s3client *S3Client) testConnection() error {

	testBucketName := fmt.Sprintf("test-%d", time.Now().UTC().UnixNano())
	b := s3client.S3Connection.Bucket(testBucketName)
	_, err := b.Get("non-existent")
	re := regexp.MustCompile("no such host")
	if re.MatchString(err.Error()) {
		return err
	}

	return nil
}
