package telegram

import (
	"log/slog"
	"time"

	"github.com/patrickmn/go-cache"
)

const CALLBACK_QUERY_COMMAND = "callbackQuery"

type bot struct {
	client          apiCaller
	cache           *cache.Cache
	cacheExparation time.Duration
	commandHandlers map[string]CommandHandler
	chatHandlers    map[string]map[int]CommandStepHandler
}

func NewBot(token string, logger *slog.Logger) *bot {
	client := newClient(token, logger)
	cache_ := cache.New(10*time.Minute, 10*time.Minute)
	commandHandlers := make(map[string]CommandHandler)
	chatHandlers := make(map[string]map[int]CommandStepHandler)

	return &bot{client, cache_, cache.DefaultExpiration, commandHandlers, chatHandlers}
}

func (bot *bot) RegisterCommandHandler(command string, handler CommandHandler) error {
	bot.commandHandlers[command] = handler

	return nil
}

func (bot *bot) RegisterCallbackQueryHandler(handler CommandHandler) error {
	bot.commandHandlers[CALLBACK_QUERY_COMMAND] = handler

	return nil
}
