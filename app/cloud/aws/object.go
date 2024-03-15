package aws

import (
	"os"
	"io"
	"fmt"
	"bytes"
	"context"
	//"strings"
	"crypto/md5"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"

	"khranity/app/utils"
	"khranity/app/log"
)

const (
	maxSize		= int64(3 * 1024 * 1024 * 1024) // 3 Gb - max size of object to upload
	maxPartSize	= int64(99 * 1024 * 1024) 		// 99 Mb- max part size for multu upload
	maxRetries  = 3								// 3 	- max count tries for upload
)


// UploadObject reads from a file and puts the data into an object in a bucket.
func (cld *Cloud) UploadObject(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	// TODO check hash sums
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

	return nil
}

// DownloadObject gets an object from a bucket and stores it in a local file.
func (cld *Cloud) DownloadObject(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	// TODO check hash sums
	try := 1
	for try <= maxRetries {
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
		hash:= md5.New()
		tee	:= io.TeeReader(result.Body, hash)
		
		size, err := io.Copy(file, tee)
		if err != nil {
			cld.logger.Debug("ioCopy: DownloadObject",
				log.String("err",	err.Error()),
				log.Int("try", 		try),
				log.Int64("size",	size),
			)
			if try == maxRetries {
				return utils.ErrExceededRetries
			} else {
				try++
			}
		} else {
			hashStr := fmt.Sprintf("%x", hash.Sum(nil))
			cld.logger.Debug("download done",
				log.String("etag", 	*result.ETag),
				log.String("hash", 	hashStr),
				log.Int64("size", 	size),
			)
			return nil // ok
		}
	}
	return nil
}

// UploadBigObject reads from a BIG file and puts the data into an object in a bucket.
func (cld *Cloud) UploadBigObject(ctx context.Context, bucketName string, objectKey string, fileName string) error {
	// TODO create multithread upload
	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	
	defer file.Close()

	// Get file size
	stats, _ := file.Stat()
	fileSize := stats.Size()
	
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

	// put file in byteArray
	buffer := make([]byte, fileSize) // wouldn't want to do this for a large file because it would store a potentially super large file into memory
	file.Read(buffer)

	output, err := cld.client.CreateMultipartUpload(context.TODO(), &s3.CreateMultipartUploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return err
	}

	parts 	:= []s3Types.CompletedPart{}
	rest 	:= fileSize
	count 	:= int32(1)

	for start := int64(0); rest != 0; start += maxPartSize {
		var size int64
		if rest < maxPartSize {
			size = rest
		} else {
			size = maxPartSize
		}
		part, err := cld.upload(output, buffer[start:start+size], count)
		if err != nil {
			_, err1 := cld.client.AbortMultipartUpload(context.TODO(), &s3.AbortMultipartUploadInput{
				Bucket:   output.Bucket,
				Key:      output.Key,
				UploadId: output.UploadId,
			})
			if err1 != nil {
				return err1
			}
			return err
		}
		rest -= size
		//fmt.Printf("Part %v complete, %v bytes remaining\n", partNum, remaining)
		parts = append(parts, part)
		count++
	}

	_, err = cld.client.CompleteMultipartUpload(context.TODO(), &s3.CompleteMultipartUploadInput{
		Bucket:   output.Bucket,
		Key:      output.Key,
		UploadId: output.UploadId,
		MultipartUpload: &s3Types.CompletedMultipartUpload{
			Parts: parts,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (cld *Cloud) upload(output *s3.CreateMultipartUploadOutput, fileBytes []byte, partNum int32) (completedPart s3Types.CompletedPart, err error) {
	try := 1
	for try <= maxRetries {
		length := int64(len(fileBytes))
		uploadResp, err := cld.client.UploadPart(context.TODO(), &s3.UploadPartInput{
			Body:          bytes.NewReader(fileBytes),
			Bucket:        aws.String(*output.Bucket),
			Key:           aws.String(*output.Key),
			PartNumber:    &partNum,
			UploadId:      output.UploadId,
			ContentLength: &length,
		})
		if err != nil {
			cld.logger.Debug("cld.client.UploadPart: upload",
				log.String("err",	err.Error()),
				log.Int("try", 		try),
			)
			if try == maxRetries {
				return s3Types.CompletedPart{}, utils.ErrExceededRetries
			} else {
				try++
			}
		} else {
			return s3Types.CompletedPart{
				ETag:       uploadResp.ETag,
				PartNumber: &partNum,
			}, nil
		}
	}
	return s3Types.CompletedPart{}, nil
}