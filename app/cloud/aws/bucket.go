package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/service/s3"

)

// BucketExists checks whether a bucket exists in the current account.
func (cld *Cloud) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	_, err := cld.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})

	if err != nil {
		return false, err
	}
	return true, err
}

// CreateBucket creates a bucket with the specified name in the specified Region.
func (cld *Cloud) CreateBucket(ctx context.Context, name string/*, region string*/) error {
	_, err := cld.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(name),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			//LocationConstraint: types.BucketLocationConstraint(region),
		},
	})
	
	return err
}