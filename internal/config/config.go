package config

import (
	"log"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env          string `yaml:"env" env:"ENV" env-default:"local"`
	HTTPServer   `yaml:"http_server" env-required:"true"`
	RouterConfig `yaml:"router_config" env-required:"true"`
	PGConfig     `yaml:"pg_config" env-required:"true"`
	HashConfig   `yaml:"hash_config" env-required:"true"`
	SmtpConfig   `yaml:"smtp_config" env-required:"true"`
	NotifyConfig `yaml:"notify_config" env-required:"true"`
	OpenAI       `yaml:"open_ai" env-required:"true"`
	Prompts      `yaml:"prompts" env-required:"true"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
}

type RouterConfig struct {
	AppPort string `yaml:"app_port" env-required:"true"`
}

type PGConfig struct {
	DBHost  string `yaml:"db_host" env-required:"true"`
	DBPort  string `yaml:"db_port" env-required:"true"`
	DBName  string `yaml:"db_name" env-required:"true"`
	DBUser  string `yaml:"db_user" env-required:"true"`
	DBPass  string `yaml:"db_pass" env-required:"true"`
	MaxConn int    `yaml:"max_conn" env-default:"32"`
}

type HashConfig struct {
	SigningKey string        `yaml:"signing_key" env-required:"true"`
	TakenTTL   time.Duration `yaml:"taken_ttl" env-default:"1h"`
}

type SmtpConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	Port     int    `yaml:"port" env-required:"true"`
	Password string `yaml:"app_password" env-required:"true"`
	Email    string `yaml:"email" env-required:"true"`
}

type NotifyConfig struct {
	Link    string `yaml:"website_link" env-required:"true"`
	Time    string `yaml:"time" env-required:"true"`
	From    string `yaml:"from" env-required:"true"`
	Subject string `yaml:"subject" env-required:"true"`
}

type OpenAI struct {
	ApiKeys []string `yaml:"api_keys" env-required:"true"`
	Ind     int32
}

type Prompts struct {
	GenerateTask     string `yaml:"generate_task" env-required:"true"`
	CheckCodeForTask string `yaml:"check_code_for_task" env-required:"true"`
	GiveHintForTask  string `yaml:"give_hint_for_task" env-required:"true"`
	Notification     string `yaml:"notification" env-required:"true"`
}

func MustLoad() *Config {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	if _, err = os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("CONFIG_PATH does not exist: %s", configPath)
	}

	var cfg Config

	if err = cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	return &cfg
}
