package aws

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"khranity/app/log"
	"khranity/app/lore"
)

var (
	once   sync.Once
	client *s3.Client
)

type Cloud struct {
	logger 	*log.Logger
	client	*s3.Client
}

func CreateClient(cloud *lore.Cloud) *s3.Client {
	resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           cloud.Url,
			SigningRegion: cloud.Region,
		}, nil
	})

	provider := credentials.NewStaticCredentialsProvider(
		lore.GetToken(cloud.AccessID),
		lore.GetToken(cloud.SecretKey),
		lore.GetToken(cloud.Token),
	)

	config, err := config.LoadDefaultConfig(
		context.TODO(), 
		config.WithEndpointResolverWithOptions(resolver), 
		config.WithCredentialsProvider(provider), 
		config.WithRegion("auto"),
	)
	if err != nil {
		return nil
	}

	return s3.NewFromConfig(config, func(o *s3.Options){
		o.UsePathStyle = true
	})
}

// NewCloud ...
func NewCloud(logger *log.Logger, cloud *lore.Cloud) *Cloud {
	once.Do(func() {
		client = CreateClient(cloud)
	})
	return &Cloud{
		logger: logger,
		client: client,
	}
}

// TestCloud ...
func TestCloud(cloud *lore.Cloud) error {
	list, err := CreateClient(cloud).ListBuckets(context.TODO(), &s3.ListBucketsInput{})
	if err != nil {
		return err
	}
	for _, bucket := range list.Buckets {
		log.Get().Debug("bucket", log.String("Name", *bucket.Name))
	}
	return nil
}
