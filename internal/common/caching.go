package common

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type ChatContext struct {
	chatId        string
	userResponses []string

	// next handler to call in conversation with bot
	nextHandler string
}

type ChatCache struct {
	cache           *cache.Cache
	cacheExparation time.Duration
}

func NewChatCache() *ChatCache {
	cache_ := cache.New(10*time.Minute, 10*time.Minute)

	return &ChatCache{cache_, cache.DefaultExpiration}
}

func newChatContext(chatId string) *ChatContext {
	return &ChatContext{chatId, []string{}, ""}
}

func (ctx *ChatContext) AppendText(userResponse string) error {
	ctx.userResponses = append(ctx.userResponses, userResponse)
	return nil
}

func (ctx *ChatContext) GetTexts() []string {
	return ctx.userResponses
}

func (ctx *ChatContext) GetNextHandler() string {
	return ctx.nextHandler
}

func (ctx *ChatContext) SetNextHandler(nextHandler string) string {
	if nextHandler == "" {
		ctx.Reset()
	}
	ctx.nextHandler = nextHandler
	return ctx.nextHandler
}

func (ctx *ChatContext) Reset() error {
	ctx.nextHandler = ""
	ctx.userResponses = []string{}
	return nil
}

func (c *ChatCache) GetOrCreateChatContext(chatId string) *ChatContext {
	ctx, found := c.cache.Get(chatId)

	if found {
		return ctx.(*ChatContext)
	}

	newCtx := newChatContext(chatId)

	c.cache.Set(chatId, newCtx, c.cacheExparation)

	return newCtx
}
