package lore

import (
	"os"
	"sync"
	"strings"

	"github.com/goccy/go-yaml"
	"go.uber.org/zap/zapcore"

	"khranity/app/log"
	"khranity/app/config"
	"khranity/app/utils"
)

var (
	lore    *Lore
	//jobs 	[]Job
	once   	sync.Once
)

type (
	Job struct {
		Name		string	`yaml:"name"`
		Path		string	`yaml:"path"`
		Schedule	string	`yaml:"schedule"`	// in cron format
		TZ			string	`yaml:"tz"`			// "Asia/Shanghai"
		//Size		int		`yaml:"size"`
		Cloud		string	`yaml:"cloud"`		// from struct Cloud.Name
		Token		string	`yaml:"token"`		// token, file:token_file
	}

	Cloud struct {
		Name		string	`yaml:"name"`
		Method		string	`yaml:"method"`			// aws, yandex, selectel
		Url			string	`yaml:"url"`
		Region		string	`yaml:"region"`
		Bucket		string	`yaml:"bucket"`
		AccessID	string	`yaml:"access_id"`		// access_key, file:access_key_file
		SecretKey	string	`yaml:"secret_key"`		// secret_key, file:secret_key_file
		Token		string	`yaml:"token"`			// token, file:token_file
	}

	Setup struct {
		OS			string	`yaml:"os"`			// nix, win
		//Cloud		string	`yaml:"cloud"`		// aws, yandex, selectel
	}

	Lore struct {
		Jobs  	[]Job		`yaml:"jobs"`
		Clouds 	[]Cloud		`yaml:"clouds"`
		Setup   Setup		`yaml:"setup"`
	}
)

func (job *Job) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Name",		job.Name)
    enc.AddString("Path",		job.Path)
	enc.AddString("Schedule",	job.Schedule)
	enc.AddString("Cloud",		job.Cloud)
	enc.AddString("Token",		job.Token)
	//enc.AddInt("Size",			job.Size)
	return nil
}

func (cloud *Cloud) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("Name",		cloud.Name)
    enc.AddString("Url",		cloud.Url)
	enc.AddString("Region",		cloud.Region)
	enc.AddString("Bucket",		cloud.Bucket)
	enc.AddString("AccessID",	cloud.AccessID)
	enc.AddString("SecretKey",	cloud.SecretKey)
	enc.AddString("Token",		cloud.Token)
	return nil
}

func Get() *Lore {
	once.Do(func() {
		lore 		= &Lore{}
		loreFile	:= config.Get().LoreFile
		data, err := os.ReadFile(loreFile)
		if err != nil {
				log.Fatal("lore.Get() error read lore config file", 
					log.String("err", 	err.Error()),
					log.String("file", 	loreFile),
				)
		}

		err = yaml.Unmarshal(data, lore)
		if err != nil {
				log.Fatal("lore.Get() error unmarshal lore config file", 
					log.String("err", 	err.Error()),
					log.String("file", 	loreFile),
				)
		}
	})
	return lore
}

func Load(loreFile string) *Lore {
	once.Do(func() {
		lore 	= &Lore{}
		data, err := os.ReadFile(loreFile)
		if err != nil {
				log.Fatal("lore.Get() error read lore-file", 
					log.String("err", 	err.Error()),
					log.String("file", 	loreFile),
				)
		}

		err = yaml.Unmarshal(data, lore)
		if err != nil {
				log.Fatal("lore.Get() error unmarshal lore-file", 
					log.String("err", 	err.Error()),
					log.String("file", 	loreFile),
				)
		}
	})
	return lore
}

func GetJobs() []Job {
	if (lore != nil){
		return lore.Jobs
	}

	return nil
}

func GetJob(name string) (*Job, error) {
	if (lore != nil) {
		for _, j := range lore.Jobs {
			if (j.Name == name) {
				return &j, nil
			}
		}
	}
	
	return nil, utils.ErrUndefinedCloud
}

func GetCloud(name string) (*Cloud, error) {
	if (lore != nil) {
		for _, c := range lore.Clouds {
			if (c.Name == name) {
				return &c, nil
			}
		}
	}
	
	return nil, utils.ErrUndefinedCloud
}

func GetToken(name string) string {
	s := strings.Split(name, "file:")
	if len(s) < 2 {
		return name
	}
	if (len(s[1]) < 1) {
		return name
	}

	data, err := os.ReadFile(strings.TrimSpace(s[1]))
	if (err != nil) {
		return name
	}

	return string(data)
}