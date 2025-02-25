package yandex

import (
	"context"
)

// UploadObject reads from a file and puts the data into an object in a bucket.
func (cld *Cloud) UploadObject(ctx context.Context, bucketName string, objectKey string, fileName string) error {

	return nil
}

// DownloadObject gets an object from a bucket and stores it in a local file.
func (cld *Cloud) DownloadObject(ctx context.Context, bucketName string, objectKey string, fileName string) error {

	return nil
}

// UploadObject reads from a BIG file and puts the data into an object in a bucket.
func (cld *Cloud) UploadBigObject(ctx context.Context, bucketName string, objectKey string, fileName string) error {

	return nil
}

// DownloadObject gets an object from a bucket and stores it in a local BIG file.
func (cld *Cloud) DownloadBigObject(ctx context.Context, bucketName string, objectKey string, fileName string) error {

	return nil
}