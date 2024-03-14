package nix

import (
	"context"
	"strings"
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron"

	"khranity/app/config"
	"khranity/app/log"
	"khranity/app/lore"
	"khranity/app/targz"
	"khranity/app/crypt"
	"khranity/app/utils"
	"khranity/app/cloud"
)

const (
	maxSize		= int64(3 * 1024 * 1024 * 1024) // 3 Gb - max size of object to upload
	maxPartSize	= int64(99 * 1024 * 1024) 		// 99 Mb- max part size for multu upload
)

func (sh *ShellNix) Start(ctx context.Context, job *lore.Job) error {
	if (len(job.Schedule) < 1) {
		sh.logger.Error("empty schedule", 
			log.Object("job", job),
		)
		return utils.ErrIncorrectJob
	}
	task := gocron.NewScheduler(time.UTC)
	
	_, err := task.Cron(job.Schedule).Do(sh.Exec, ctx, job)
	if err != nil {
		log.Debug("shell.Start", log.String("err", err.Error()))
		return utils.ErrIncorrectJob
	}
	
	task.StartAsync()

	return nil
}

func (sh *ShellNix) Exec(ctx context.Context, job *lore.Job) error {
	temp := config.Get().TempFolder
	if len(temp) < 1 {
		temp = os.TempDir()
	}

	err := sh.Put(ctx, job, temp)
	if (err != nil) {
		sh.logger.Error("job failed",
			log.String("err", 	err.Error()),
			log.Object("job", 	job),
		)
		return err
	}

	sh.logger.Info("job done",
		log.Object("job", 	job),
	)
	
	return nil
}

func (sh *ShellNix) Get(ctx context.Context, job *lore.Job, temp string) error {
	start 	:= time.Now()
	// create tar-file
	path := strings.Replace(job.Path, " ", "_", -1)
	name := strings.Replace(job.Name, " ", "_", -1)
	isDir, err := isDir(path)
	if (err != nil) || !isDir {
		sh.logger.Error("incorrect path: shell.Get", 
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("path", 		path),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrIncorrectJob
	}

	fileName := fmt.Sprintf("%v/%v.tar.gz", temp, name)
	
	// download from cloud
	jobCloud, err := lore.GetCloud(job.Cloud)
	if (err != nil) {
		sh.logger.Error("lore.GetCloud",
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrUndefinedCloud
	}

	cld, err := cloud.New(ctx, sh.logger, jobCloud)
	if (err != nil) {
		sh.logger.Error("cloud.New",
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrCloudInternal
	}

	err = cld.Cloud.DownloadObject(ctx, jobCloud.Bucket, job.Name, fileName)
	if (err != nil) {
		sh.logger.Error("cloud.DownloadObject",
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("file", 		fileName),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrCloudInternal
	}

	// decrypt if necessary
	token := lore.GetToken(job.Token)
	//if (len(token) > 15) {
	err = crypt.DecryptFile(fileName, fileName, token)
	if (err != nil) {
		sh.logger.Error("decrypt failed", 
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("token",		token),
			log.String("file", 		fileName),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrInternal
	}
	//}

	// extract archive
	err = targz.Extract(fileName, path)
	if err != nil {
		sh.logger.Error("extract failed", 
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("path", 		path),
			log.String("file", 		fileName),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrInternal
	}

	// remove temp-file
	err = os.Remove(fileName)
    if err != nil {
		sh.logger.Error("os.Remove",
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("file", 		fileName),
			log.String("duration", 	time.Since(start).String()),
		)
    }

	sh.logger.Info("job done",
		log.String("method", 	"Get"),
		log.Object("job", 		job),
		log.String("duration", 	time.Since(start).String()),
	)

	return nil
}

func (sh *ShellNix) Put(ctx context.Context, job *lore.Job, temp string) error {
	start 	:= time.Now()
	// create tar-file
	path := strings.Replace(job.Path, " ", "_", -1)
	name := strings.Replace(job.Name, " ", "_", -1)
	isDir, err := isDir(path)
	if (err != nil) || !isDir {
		sh.logger.Error("incorrect path: shell.Put", 
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("path", 		path),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrIncorrectJob
	}

	// tz, err := time.LoadLocation(job.TZ)
	// if err != nil { 
	// 	tz = time.UTC
	// }
	
	// tail := time.Now().In(tz).Format("2006_01_02__15_04_05")
	// if (len(temp) < 1) {
	// 	temp = "temp"
	// }
	//fileName := fmt.Sprintf("%v/%v_%v.tar.gz", temp, name, tail)
	fileName := fmt.Sprintf("%v/%v.tar.gz", temp, name)
	
	err = targz.Compress(path, fileName)
	if err  != nil {
		sh.logger.Error("compress failed: shell.Put", 
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("file", 		fileName),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrInternal
	}

	// TODO files more then 99 MB skip
	stat, err := os.Stat(fileName)
	if (err != nil) {
		sh.logger.Error("fileinfo failed: shell.Put", 
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("file", 		fileName),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrInternal
	}
	
	// encrypt if necessary
	token := lore.GetToken(job.Token)
	//if (len(token) > 15) {
	err = crypt.EncryptFile(fileName, fileName, token)
	if (err != nil) {
		sh.logger.Error("encrypt failed: shell.Put", 
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("token",		token),
			log.String("file", 		fileName),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrInternal
	}
	//}

	// upload to cloud
	jobCloud, err := lore.GetCloud(job.Cloud)
	if (err != nil) {
		sh.logger.Error("lore.GetCloud: shell.Put",
			log.String("err", 	err.Error()),
			log.Object("job", 	job),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrUndefinedCloud
	}

	cld, err := cloud.New(ctx, sh.logger, jobCloud)
	if (err != nil) {
		sh.logger.Error("cloud.New: shell.Put",
			log.String("err", 	err.Error()),
			log.Object("job", 	job),
			log.String("duration", 	time.Since(start).String()),
		)
		return utils.ErrCloudInternal
	}

	if stat.Size() > maxSize {
		return utils.ErrOverSize
	}

	if stat.Size() > maxPartSize {
		err = cld.Cloud.UploadBigObject(ctx, jobCloud.Bucket, job.Name, fileName)
		if (err != nil) {
			sh.logger.Error("cloud.UploadBigObject: shell.Put",
				log.String("err", 		err.Error()),
				log.Object("job", 		job),
				log.String("file", 		fileName),
				log.String("duration", 	time.Since(start).String()),
			)
			return utils.ErrCloudInternal
		}
	} else {
		err = cld.Cloud.UploadObject(ctx, jobCloud.Bucket, job.Name, fileName)
		if (err != nil) {
			sh.logger.Error("cloud.UploadObject: shell.Put",
				log.String("err", 		err.Error()),
				log.Object("job", 		job),
				log.String("file", 		fileName),
				log.String("duration", 	time.Since(start).String()),
			)
			return utils.ErrCloudInternal
		}
	}

	// remove temp-file
	err = os.Remove(fileName)
    if err != nil {
		sh.logger.Error("os.Remove: shell.Put",
			log.String("err", 		err.Error()),
			log.Object("job", 		job),
			log.String("file", 		fileName),
			log.String("duration", 	time.Since(start).String()),
		)
    }

	sh.logger.Info("job done",
		log.String("method", 	"Put"),
		log.Object("job", 		job),
		log.String("duration", 	time.Since(start).String()),
	)
	
	return nil
}

func isDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	return fileInfo.IsDir(), err
}