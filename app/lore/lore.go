package lore

import (
	"io"
	"os"
	"sync"
	"time"
	"bytes"
	"strings"
	"crypto/sha256"

	"github.com/goccy/go-yaml"
	"go.uber.org/zap/zapcore"

	"khranity/app/config"
	"khranity/app/log"
	"khranity/app/utils"
)

type State struct {
	cond *sync.Cond
	hash []byte
}

var (
	lore  *Lore
	once  sync.Once
	state *State
)

type Ignores []string

type (
	Job struct {
		Name     string `yaml:"name"`
		Path     string `yaml:"path"`
		Schedule string `yaml:"schedule"` // in cron format
		TZ       string `yaml:"tz"`       // "Asia/Shanghai"
		//Size		int			`yaml:"size"`
		Cloud    string  `yaml:"cloud"`    // from struct Cloud.Name
		Archiver string  `yaml:"archiver"` // archiver method
		Token    string  `yaml:"token"`    // token, file:token_file
		Ignore   Ignores `yaml:"ignore"`   // ignore object list
	}

	Cloud struct {
		Name      string `yaml:"name"`
		Method    string `yaml:"method"` // aws, yandex, selectel
		Url       string `yaml:"url"`
		Region    string `yaml:"region"`
		Bucket    string `yaml:"bucket"`
		AccessID  string `yaml:"access_id"`  // access_key, file:access_key_file
		SecretKey string `yaml:"secret_key"` // secret_key, file:secret_key_file
		Token     string `yaml:"token"`      // token, file:token_file
	}

	Setup struct {
		OS string `yaml:"os"` // nix, win
	}

	Lore struct {
		Jobs   []Job   `yaml:"jobs"`
		Clouds []Cloud `yaml:"clouds"`
		Setup  Setup   `yaml:"setup"`
	}
)

func (data Ignores) MarshalLogArray(arr zapcore.ArrayEncoder) error {
	if data == nil {
		return utils.ErrInternal
	}
	for i := range data {
		arr.AppendString(data[i])
	}
	return nil
}

func (job *Job) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if job == nil {
		return utils.ErrIncorrectJob
	}
	enc.AddString("Name", job.Name)
	enc.AddString("Path", job.Path)
	enc.AddArray("Ignore", job.Ignore)
	enc.AddString("Schedule", job.Schedule)
	enc.AddString("Cloud", job.Cloud)
	enc.AddString("Token", job.Token)
	enc.AddString("Archiver", job.Archiver)
	//enc.AddInt("Size",			job.Size)
	return nil
}

func (cloud *Cloud) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if cloud == nil {
		return utils.ErrUndefinedCloud
	}
	enc.AddString("Name", cloud.Name)
	enc.AddString("Url", cloud.Url)
	enc.AddString("Region", cloud.Region)
	enc.AddString("Bucket", cloud.Bucket)
	enc.AddString("AccessID", cloud.AccessID)
	enc.AddString("SecretKey", cloud.SecretKey)
	enc.AddString("Token", cloud.Token)
	return nil
}

func (lore *Lore) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	if lore == nil {
		return utils.ErrUndefinedCloud
	}
	for _, cloud := range lore.Clouds {
		enc.AddObject("Cloud", &cloud)
	}
	return nil
}

func (c *State) Wait() {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()
	for {
		c.cond.Wait()
		load()
	}
}

func (c *State) Broadcast() {
	c.cond.L.Lock()
	defer c.cond.L.Unlock()
	c.cond.Broadcast()
}

func needUpdating() bool {
	f, err := os.Open(config.Get().LoreFile)
	if err != nil {
		log.Error("needUpdating(): read lore file failed",
			log.String("err", err.Error()),
		)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Error("needUpdating(): copy lore file failed",
			log.String("err", err.Error()),
		)
	}

	if !bytes.Equal(state.hash, h.Sum(nil)) {
		state.hash = h.Sum(nil)
		return true
	}

	return false
}

func load() {
	loreFile := config.Get().LoreFile
	data, err := os.ReadFile(loreFile)
	if err != nil {
		log.Fatal("lore.Get() read lore file failed",
			log.String("err", err.Error()),
			log.String("file", loreFile),
		)
	}

	err = yaml.Unmarshal(data, lore)
	if err != nil {
		log.Fatal("lore.Get() unmarshal lore file failed",
			log.String("err", err.Error()),
			log.String("file", loreFile),
		)
	}
	log.Debug("lore updated",
		log.Object("lore", lore),
	)
}

func CheckUpdating() {
	if needUpdating() {
		state.Broadcast()
	}
}

func Get() *Lore {
	once.Do(func() {
		lore	= &Lore{}
		state = &State{}
		if needUpdating() {
			load()
		}
    	state.cond = sync.NewCond(&sync.Mutex{})
		go func() {
			state.Wait()
		}()
		time.Sleep(1 * time.Second)
	})
	
	return lore
}

func Load(loreFile string) *Lore {
	once.Do(func() {
		lore = &Lore{}
		data, err := os.ReadFile(loreFile)
		if err != nil {
			log.Fatal("lore.Get() error read lore-file",
				log.String("err", err.Error()),
				log.String("file", loreFile),
			)
		}

		err = yaml.Unmarshal(data, lore)
		if err != nil {
			log.Fatal("lore.Get() error unmarshal lore-file",
				log.String("err", err.Error()),
				log.String("file", loreFile),
			)
		}
	})
	return lore
}

func GetJobs() []Job {
	if lore != nil {
		return lore.Jobs
	}

	return nil
}

func GetJob(name string) (*Job, error) {
	if lore != nil {
		for _, j := range lore.Jobs {
			if j.Name == name {
				return &j, nil
			}
		}
	}

	return nil, utils.ErrUndefinedCloud
}

func GetCloud(name string) (*Cloud, error) {
	if lore != nil {
		for _, c := range lore.Clouds {
			if c.Name == name {
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
	if len(s[1]) < 1 {
		return name
	}

	data, err := os.ReadFile(strings.TrimSpace(s[1]))
	if err != nil {
		return name
	}

	return string(data)
}
