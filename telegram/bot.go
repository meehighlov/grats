package telegram

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type bot struct {
	client          apiCaller
	cache           *cache.Cache
	cacheExparation time.Duration
	commandHandlers map[string]CommandHandler
	chatHandlers    map[string]map[int]CommandStepHandler
}

func NewBot(token string) *bot {
	client := newClient(token, nil)
	cache_ := cache.New(10*time.Minute, 10*time.Minute)
	commandHandlers := make(map[string]CommandHandler)
	chatHandlers := make(map[string]map[int]CommandStepHandler)

	return &bot{client, cache_, cache.DefaultExpiration, commandHandlers, chatHandlers}
}

func (bot *bot) RegisterCommandHandler(command string, handler CommandHandler) error {
	bot.commandHandlers[command] = handler

	return nil
}
