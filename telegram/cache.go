package telegram

import (
	"time"

	"github.com/patrickmn/go-cache"
)

type BotCache struct {
	cache *cache.Cache
}

func NewBotCache() BotCacheAccess {
	return &BotCache{cache: cache.New(10*time.Minute, 10*time.Minute)}
}

type BotCacheAccess interface {
	GetOrCreateChatContext(string) *context
}

func (botCache *BotCache) GetOrCreateChatContext(chatId string) *context {
	ctx, found := botCache.cache.Get(chatId)

	if found {
		return ctx.(*context)
	}

	newCtx := newChatContext(chatId)

	botCache.cache.Set(chatId, newCtx, cache.DefaultExpiration)

	return newCtx
}

type ChatContext interface {
	AppendUserResponse(string) error
	GetUserResponses() []string
	GetCommandInProgress() *string
	SetCommandInProgress(string) error
	GetStepDone() int
	SetStepDone(int) error
	Reset() error
}

type context struct {
	chatId string
	userResponses []string
	commandInProgress *string
	stepDone int
}

func newChatContext(chatId string) *context {
	return &context{chatId, []string{}, nil, 0}
}

func (ctx *context) AppendUserResponse(userResponse string) error {
	ctx.userResponses = append(ctx.userResponses, userResponse)
	return nil
}

func (ctx *context) GetUserResponses() []string {
	return ctx.userResponses
}

func (ctx *context) GetCommandInProgress() *string {
	return ctx.commandInProgress
}

func (ctx *context) GetStepDone() int {
	return ctx.stepDone
}

func (ctx *context) SetStepDone(stepDone int) error {
	ctx.stepDone = stepDone
	return nil
}

func (ctx *context) SetCommandInProgress(command string) error {
	if ctx.commandInProgress != nil {
		if *ctx.commandInProgress != command {
			ctx.Reset()
		}
	}
	ctx.commandInProgress = &command
	return nil
}

func (ctx *context) Reset() error {
	ctx.commandInProgress = nil
	ctx.userResponses = []string{}
	ctx.stepDone = 0
	return nil
}
