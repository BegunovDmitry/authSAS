package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	AppMode         string            `yaml:"app_mode" env-required:"true"`
	PermStoragePath string            `yaml:"permanent_storage_path" env-required:"true"`
	JWTTokenTTL     time.Duration     `yaml:"jwt_token_ttl" env-default:"24h"`
	JWTSecret     string    `yaml:"jwt_secret" env-required:"true"`
	Grpc            GrpcCnofig        `yaml:"grpc"`
	TempStorage     TempStorageConfig `yaml:"temp_storage"`
	EmailSender EmailSender `yaml:"email_sender"`
}

type GrpcCnofig struct {
	Domain         string        `yaml:"domain" env-required:"true"`
	Port           int           `yaml:"port" env-required:"true"`
	RequestTimeout time.Duration `yaml:"req_timeout" env-default:"1m"`
}

type TempStorageConfig struct {
	TempStoragePath string `yaml:"temporary_storage_path" env-required:"true"`
	CodeTTL  time.Duration `yaml:"code_ttl" env-default:"10m"`
}

type EmailSender struct {
	Email string `yaml:"email" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
}

func MustLoad() *Config {
	path := fillConfigPath()

	if path == "" {
		panic("Config file path is empty")
	}

	if _, err := os.Stat(path); err != nil{
		panic("Config file is does not exist: "+path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("Failed to read config file: "+err.Error())
	}

	return &cfg
}

func fillConfigPath() string {
	var path string
	
	flag.StringVar(&path, "config", "", "Path to config .yaml file")
	flag.Parse()

	if path == "" {
		path = os.Getenv("CONFIG_PATH")
	}

	return path
}