package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

const configPath = "/auth_service/config/local.yaml"

type (
	Config struct {
		Env string `yaml:"env" env-default:"local"`
		UsersDB
		HTTPServer `yaml:"http_server"`
	}
	HTTPServer struct {
		Port string `yaml:"address" env-default:"8080"`
	}
	UsersDB struct {
		DataSourceName string `env:"DB_SOURCE_NAME" env-required:"true"`
		Name           string `yaml:"name" env-default:"usersdb"`
	}
)

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("MustLoad: unable to load cfg: %s", err)
	}
	return &cfg
}
