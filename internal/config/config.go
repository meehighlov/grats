package config

import (
	"log"
	"os"

	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	ENV                   string `env:"ENV" env-default:"local"`
	BotToken              string `env:"BOT_TOKEN" env-required:"true"`
	BotName               string `env:"BOT_NAME" env-required:"true"`
	Admins                string `env:"ADMINS" env-required:"true"`
	ReportChatId          string `env:"REPORT_CHAT_ID" env-required:"true"`
	HandlerExecTimeoutSec int    `env:"HANDLER_EXEC_TIMEOUT_SEC" env-default:"2"`
	Timezone              string `env:"TIMEZONE" env-default:"Europe/Moscow"`
	SupportChatId         string `env:"SUPPORT_CHAT_ID" env-required:"true"`

	UseWebhook         bool   `env:"USE_WEBHOOK" env-default:"false"`
	WebhookAddr        string `env:"WEBHOOK_ADDR" env-default:":8080"`
	WebhookSecretToken string `env:"WEBHOOK_SECRET_TOKEN" env-default:""`
	WebhookTlsOn       bool   `env:"WEBHOOK_TLS_ON" env-default:"false"`
	WebhookTlsCertFile string `env:"WEBHOOK_TLS_CERT_FILE" env-default:""`
	WebhookTlsKeyFile  string `env:"WEBHOOK_TLS_KEY_FILE" env-default:""`
	WebhookTlsAddr     string `env:"WEBHOOK_TLS_ADDR" env-default:":443"`

	loaded bool `env-default:"false"`
}

func (cfg *Config) AdminList() []string {
	return strings.Split(cfg.Admins, ",")
}

func (cfg *Config) HandlerTmeout() time.Duration {
	return time.Duration(cfg.HandlerExecTimeoutSec) * time.Second
}

var cfg Config

// loads config from .env
// also sets TZ env variable from according .env value
func MustLoad() *Config {
	if _, err := os.Stat("env/grats/.env"); os.IsNotExist(err) {
		log.Fatal("Not found .env file")
	}

	err := cleanenv.ReadConfig("env/grats/.env", &cfg)
	if err != nil {
		log.Fatal("Failed to read envs:", err.Error())
	}

	os.Setenv("TZ", cfg.Timezone)

	cfg.loaded = true

	return &cfg
}

func Cfg() *Config {
	if !cfg.loaded {
		log.Fatal("Accessing not loaded config. Exiting.")
	}

	return &cfg
}
