package common

import (
	"time"

	"github.com/patrickmn/go-cache"
)

const (
	ENTRYPOINT_STEP = "1"
	STEPS_DONE = "done"
)

type ChatContext struct {
	chatId            string
	userResponses     []string
	commandInProgress string
	stepTODO          string
}

type ChatCache struct {
	cache *cache.Cache
	cacheExparation time.Duration
}

func NewChatCache() *ChatCache {
	cache_ := cache.New(10*time.Minute, 10*time.Minute)

	return &ChatCache{cache_, cache.DefaultExpiration}
}

func newChatContext(chatId string) *ChatContext {
	return &ChatContext{chatId, []string{}, "", ENTRYPOINT_STEP}
}

func (ctx *ChatContext) AppendText(userResponse string) error {
	ctx.userResponses = append(ctx.userResponses, userResponse)
	return nil
}

func (ctx *ChatContext) GetTexts() []string {
	return ctx.userResponses
}

func (ctx *ChatContext) GetCommandInProgress() string {
	return ctx.commandInProgress
}

func (ctx *ChatContext) GetStepTODO() string {
	return ctx.stepTODO
}

func (ctx *ChatContext) SetStepTODO(step string) error {
	ctx.stepTODO = step
	return nil
}

func (ctx *ChatContext) SetCommandInProgress(command string) error {
	if ctx.commandInProgress != "" {
		if ctx.commandInProgress != command {
			ctx.Reset()
		}
	}
	ctx.commandInProgress = command
	return nil
}

func (ctx *ChatContext) Reset() error {
	ctx.commandInProgress = ""
	ctx.userResponses = []string{}
	ctx.stepTODO = ENTRYPOINT_STEP
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
