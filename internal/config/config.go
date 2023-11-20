package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

func MustLoad(ConfigPath string) Config { //Must это приставка которая кидает панику, типа мы только запускаем приложение поэтому мб норм уронить его
	if _, err := os.Stat(ConfigPath); os.IsNotExist(err) {
		log.Fatalf("config file %s does not exist", ConfigPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(ConfigPath, &cfg); err != nil {
		log.Fatalf("cannot read config file: %s", err)
	}

	return cfg
}
