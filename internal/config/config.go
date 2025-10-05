package config

import (
	"log"
	"os"

	"strings"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type txContextKey string

const (
	PROD  = "prod"
	LOCAL = "local"
)

type Config struct {
	ENV                   string       `env:"ENV" env-default:"local"`
	PGDSN                 string       `env:"PG_DSN" env-required:"true"`
	BotName               string       `env:"BOT_NAME" env-required:"true"`
	Admins                string       `env:"ADMINS" env-required:"true"`
	ReportChatId          string       `env:"REPORT_CHAT_ID" env-required:"true"`
	HandlerExecTimeoutSec int          `env:"HANDLER_EXEC_TIMEOUT_SEC" env-default:"2"`
	Timezone              string       `env:"TIMEZONE" env-default:"Europe/Moscow"`
	SupportChatId         string       `env:"SUPPORT_CHAT_ID" env-required:"true"`
	TxKey                 txContextKey `env:"TX_KEY" env-default:"tx"`

	TelegramToken              string `env:"TELEGRAM_TOKEN" env-required:"true"`
	TelegramUseWebook          bool   `env:"TELEGRAM_USE_WEBHOOK" env-default:"false"`
	TelegramWebhookToken       string `env:"TELEGRAM_WEBHOOK_TOKEN" env-default:""`
	TelegramWebhookAddress     string `env:"TELEGRAM_WEBHOOK_ADDRESS" env-default:":8080"`
	TelegramWebhookTLSAddress  string `env:"TELEGRAM_WEBHOOK_TLS_ADDRESS" env-default:":443"`
	TelegramWebhookTLSCertFile string `env:"TELEGRAM_WEBHOOK_TLS_CERT_FILE" env-default:""`
	TelegramWebhookTLSKeyFile  string `env:"TELEGRAM_WEBHOOK_TLS_KEY_FILE" env-default:""`
	TelegramUseTLS             bool   `env:"TELEGRAM_USE_TLS" env-default:"false"`
	TelegramHandlerTimeoutSec  int    `env:"TELEGRAM_HANDLER_TIMEOUT_SEC" env-default:"2"`
	TelegramPollingWorkers     int    `env:"TELEGRAM_POLLING_WORKERS" env-default:"10"`

	RunMigrations bool   `env:"RUN_MIGRATIONS" env-default:"false"`
	MigrationsDir string `env:"MIGRATIONS_DIR" env-default:"grats-migrations"`
	ShortIDLength int    `env:"SHORT_ID_LENGTH" env-default:"6"`

	RedisAddr     string `env:"REDIS_ADDR" env-default:"localhost:6379"`
	RedisPassword string `env:"REDIS_PASSWORD" env-default:""`
	RedisDB       int    `env:"REDIS_DB" env-default:"1"`

	LoggingFileName string `env:"LOGGING_FILE_NAME" env-default:"grats.log"`

	ChatCacheExpirationMinutes int `env:"CHAT_CACHE_EXPIRATION_MINUTES" env-default:"10"`

	ListLimitLen int `env:"LIST_LIMIT_LEN" env-default:"5"`
}

func (cfg *Config) AdminList() []string {
	return strings.Split(cfg.Admins, ",")
}

func (cfg *Config) HandlerTmeout() time.Duration {
	return time.Duration(cfg.HandlerExecTimeoutSec) * time.Second
}

func (cfg *Config) IsProd() bool {
	return cfg.ENV == PROD
}

// loads config from .env
// also sets TZ env variable from according .env value
func MustLoad() *Config {
	if _, err := os.Stat("env/grats/.env"); os.IsNotExist(err) {
		log.Fatal("Not found .env file")
	}

	var cfg Config

	err := cleanenv.ReadConfig("env/grats/.env", &cfg)
	if err != nil {
		log.Fatal("Failed to read envs:", err.Error())
	}

	os.Setenv("TZ", cfg.Timezone)

	return &cfg
}
