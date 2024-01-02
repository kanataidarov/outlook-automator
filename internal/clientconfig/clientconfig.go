package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type ClientConfig struct {
	OutlookClient `yaml:"outlook_client"`
}

type OutlookClient struct {
	AppGraphId string `yaml:"app_graph_id"`
	ClientId   string `yaml:"client_id"`
	TenantId   string `yaml:"tenant_id"`
	Uname      string `yaml:"uname"`
	Paswd      string `yaml:"paswd"`
}

func Load() *ClientConfig {
	// TODO rewrite for specifying default CONFIG_PATH
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Couldn't find Config file at %s", configPath)
	}

	var cfg ClientConfig
	// TODO replace cleanenv with more popular lib
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Couldn't read config: %s", err)
	}

	return &cfg
}
