package config

import (
	"encoding/json"
	"os"
	"log"
	"sync"

	"github.com/kelseyhightower/envconfig"
)

// Config is a config :)
type Config struct {
	LogLevel        string 	`envconfig:"LOG_LEVEL"`
	LogFolder		string 	`envconfig:"LOG_FOLDER"`
	LogPrefix       string 	`envconfig:"LOG_PREFIX"`
	LoreFile		string 	`envconfig:"LORE_FILE"`
	TempFolder		string 	`envconfig:"TEMP_FOLDER"`
}

var (
	config Config
	once   sync.Once
)

// Get reads config from environment. Once.
func Get() *Config {
	once.Do(func() {
		err := envconfig.Process("", &config)
		if err != nil {
			log.Fatal(err)
		}
		configBytes, err := json.MarshalIndent(config, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		infoLog := log.New(os.Stdout, "INIT\t", log.Ldate|log.Ltime|log.LUTC)
		infoLog.Println("Configuration:", string(configBytes))
	})
	return &config
}