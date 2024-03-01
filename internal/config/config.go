package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

type Config struct {
	Env         string `yaml:"env" env-default:"local" env-required:"true"`
	StoragePath `yaml:"storage_path"`
	HTTPServer  `yaml:"http_server"`
}

type StoragePath struct {
	DbType         string `yaml:"db_type" env-default:"postgres"`
	DbUser         string `yaml:"db_user" env-default:"userL0"`
	DbUserPassword string `yaml:"db_user_password" env-default:"123456"`
	DbPort         string `yaml:"db_port" env-default:"0.0.0.0:5432"`
	DbName         string `yaml:"db_name" env-default:"postgres"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8085"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

func MustLoad() *Config {
	// получаем путь до env файла из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")
	//configPath := "../../config/local.yaml"
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}
	//проверка существования файла конфига
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file doesn't exist: %s", configPath)
	}

	var cfg Config

	//читаем конфиг файл и заполняем структуру
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
