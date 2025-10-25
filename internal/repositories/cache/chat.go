package cache

import (
	"context"
	"encoding/json"

	"github.com/meehighlov/grats/internal/repositories/models"
)

func newChatContext(chatId string) *models.ChatContext {
	return &models.ChatContext{
		ChatId:        chatId,
		UserResponses: []string{},
		StateStatus:   "",
	}
}

// creates chat context
// if chat context not exists - creates new one
// else - return existed
func (r *Repository) createChatContext(ctx context.Context, chatId string) (*models.ChatContext, error) {
	val, err := r.redis.Redis.Get(ctx, chatId).Result()

	if err == nil {
		var ctx models.ChatContext
		if err := json.Unmarshal([]byte(val), &ctx); err == nil {
			return &ctx, nil
		}
	}

	newCtx := newChatContext(chatId)

	jsonCtx, _ := json.Marshal(newCtx)
	cmd := r.redis.Redis.Set(ctx, chatId, jsonCtx, r.redis.CacheExpiration)
	_, err = cmd.Result()
	if err != nil {
		return &models.ChatContext{}, err
	}

	return newCtx, nil
}

func (r *Repository) saveChatContext(ctx context.Context, chatContext *models.ChatContext) error {
	jsonCtx, err := json.Marshal(chatContext)
	if err != nil {
		return err
	}

	return r.redis.Redis.Set(ctx, chatContext.ChatId, jsonCtx, r.redis.CacheExpiration).Err()
}

func (r *Repository) AppendText(ctx context.Context, chatId string, text string) error {
	chatContext, err := r.createChatContext(ctx, chatId)
	if err != nil {
		return err
	}

	chatContext.AppendText(text)
	return r.saveChatContext(ctx, chatContext)
}

func (r *Repository) GetTexts(ctx context.Context, chatId string) ([]string, error) {
	chatContext, err := r.createChatContext(ctx, chatId)
	if err != nil {
		return nil, err
	}

	return chatContext.GetTexts(), nil
}

func (c *Repository) Reset(ctx context.Context, chatId string) error {
	chatContext, err := c.createChatContext(ctx, chatId)
	if err != nil {
		return err
	}

	chatContext.Reset()
	return c.saveChatContext(ctx, chatContext)
}
