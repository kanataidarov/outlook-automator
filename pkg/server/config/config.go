package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type ServerConfig struct {
	HttpServer `yaml:"http_server"`
	LogLevel   string `yaml:"log_level" env-default:"info"`
}

type HttpServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8091"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func Load() *ServerConfig {
	// TODO rewrite for specifying default CONFIG_PATH
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("Couldn't find Config file at %s", configPath)
	}

	var cfg ServerConfig
	// TODO replace cleanenv with more popular lib
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("Couldn't read config: %s", err)
	}

	return &cfg
}
