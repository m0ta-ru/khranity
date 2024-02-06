package aws

import (
	"context"
	"os"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"khranity/app/utils"
)

// UploadObject reads from a file and puts the data into an object in a bucket.
func (cld *Cloud) UploadObject(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	
	defer file.Close()
	
	if (cld.client == nil) {
		return utils.ErrInternal
	}

	exists, err := cld.BucketExists(ctx, bucketName)
	if (err != nil) || (!exists) {
		err := cld.CreateBucket(ctx, bucketName)
		if (err != nil) {
			return err
		}
	}

	_, err = cld.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		return err
	}

	return err
}

// DownloadObject gets an object from a bucket and stores it in a local file.
func (cld *Cloud) DownloadObject(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	result, err := cld.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}
	defer result.Body.Close()
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	body, err := io.ReadAll(result.Body)
	if err != nil {
		return err
	}
	_, err = file.Write(body)

	return err
}