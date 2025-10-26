package config

type Config struct {
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
	TelegramListLimitLen       int    `env:"TELEGRAM_LIST_LIMIT_LEN" env-default:"5"`
}
