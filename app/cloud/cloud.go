package cloud

import (
	"context"

	"khranity/app/cloud/aws"
	"khranity/app/cloud/selectel"
	"khranity/app/cloud/yandex"
	"khranity/app/log"
	"khranity/app/lore"
	"khranity/app/utils"
)

type Cloud interface {
	//Start(context.Context, lore.Job) error
	//Exec(context.Context, lore.Job) error
	//Upload(context.Context, string) error
	//Download()
	UploadObject(context.Context, string, string, string) error
	DownloadObject(context.Context, string, string, string) error

	UploadBigObject(context.Context, string, string, string) error
	//List(context.Context, lore.Job) error
}

// OS ...
type OS struct {
	Cloud Cloud
}

var os OS

func New(ctx context.Context, logger *log.Logger, cloud *lore.Cloud) (*OS, error) {
	switch cloud.Method {
	case "aws", "s3", "amazon":
		os.Cloud = aws.NewCloud(logger, cloud)
	case "yandex", "ya":
		os.Cloud = yandex.NewCloud(logger)
	case "sel", "selectel":
		os.Cloud = selectel.NewCloud(logger)
	default:
		return nil, utils.ErrUndefinedCloud
	}

	return &os, nil
}

func TestClouds(ctx context.Context, clouds []lore.Cloud) error {
	for _, cloud := range clouds {
		switch cloud.Method {
		case "aws", "s3", "amazon":
			err := aws.TestCloud(&cloud)
			if err != nil {
				log.Get().Error("failed testing cloud",
					log.String("err", err.Error()),
					log.Object("cloud", &cloud),
				)
				return utils.ErrCloudInternal
			}
		case "yandex", "ya":
			err := yandex.TestCloud(&cloud)
			if err != nil {
				log.Get().Error("failed testing cloud",
					log.String("err", err.Error()),
					log.Object("cloud", &cloud),
				)
				return utils.ErrCloudInternal
			}
		case "sel", "selectel":
			err := selectel.TestCloud(&cloud)
			if err != nil {
				log.Get().Error("failed testing cloud",
					log.String("err", err.Error()),
					log.Object("cloud", &cloud),
				)
				return utils.ErrCloudInternal
			}
		default:
			return utils.ErrUndefinedCloud
		}
		log.Get().Info("cloud test passed with successful",
			log.String("Name", cloud.Name),
		)
	}

	return nil
}