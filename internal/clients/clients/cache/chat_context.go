package cache

import (
	"context"
	"encoding/json"
	"errors"
)

type WorkoutHolder interface {
	GetID() string
	GetDrills() []string
}

type Workout struct {
	ID     string   `json:"id"`
	Drills []string `json:"drills"`
}

type chatContext struct {
	ChatId        string   `json:"chatId"`
	UserResponses []string `json:"userResponses"`
	Workout       *Workout `json:"workout"`
	NextHandler   string   `json:"nextHandler"`
}

func newchatContext(chatId string) *chatContext {
	return &chatContext{
		ChatId:        chatId,
		UserResponses: []string{},
		NextHandler:   "",
		Workout: &Workout{
			ID:     "",
			Drills: []string{},
		},
	}
}

func (ctx *chatContext) appendText(userResponse string) error {
	ctx.UserResponses = append(ctx.UserResponses, userResponse)
	return nil
}

func (ctx *chatContext) getTexts() []string {
	return ctx.UserResponses
}

func (ctx *chatContext) getNextHandler() string {
	return ctx.NextHandler
}

func (ctx *chatContext) reset() error {
	ctx.NextHandler = ""
	ctx.UserResponses = []string{}
	ctx.Workout = &Workout{
		ID:     "",
		Drills: []string{},
	}
	return nil
}

func (c *Client) createChatContext(ctx context.Context, chatId string) (*chatContext, error) {
	val, err := c.Redis.Get(ctx, chatId).Result()

	if err == nil {
		var ctx chatContext
		if err := json.Unmarshal([]byte(val), &ctx); err == nil {
			return &ctx, nil
		}
	}

	newCtx := newchatContext(chatId)

	jsonCtx, _ := json.Marshal(newCtx)
	cmd := c.Redis.Set(ctx, chatId, jsonCtx, c.CacheExpiration)
	_, err = cmd.Result()
	if err != nil {
		return &chatContext{}, err
	}

	return newCtx, nil
}

// creates chat context
// if chat context not exists - creates new one
// else - return existed
func (c *Client) CreateChatContext(ctx context.Context, chatId string) error {
	chatContext, err := c.createChatContext(ctx, chatId)
	if chatContext.ChatId == "" {
		return err
	}

	return nil
}

func (c *Client) GetNextHandler(ctx context.Context, chatId string) (string, error) {
	chatContext, err := c.createChatContext(ctx, chatId)
	if err != nil {
		return "", err
	}

	return chatContext.getNextHandler(), nil
}

func (c *Client) savechatContext(ctx context.Context, chatContext *chatContext) error {
	jsonCtx, err := json.Marshal(chatContext)
	if err != nil {
		return err
	}

	return c.Redis.Set(ctx, chatContext.ChatId, jsonCtx, c.CacheExpiration).Err()
}

func (c *Client) SetNextHandler(ctx context.Context, chatId string, nextHandler string) error {
	if nextHandler == "" {
		return c.Reset(ctx, chatId)
	}

	chatContext, err := c.createChatContext(ctx, chatId)
	if err != nil {
		return errors.New("chat context not found")
	}

	chatContext.NextHandler = nextHandler
	return c.savechatContext(ctx, chatContext)
}

func (c *Client) AppendText(ctx context.Context, chatId string, text string) error {
	chatContext, err := c.createChatContext(ctx, chatId)
	if err != nil {
		return err
	}

	chatContext.appendText(text)
	return c.savechatContext(ctx, chatContext)
}

func (c *Client) GetTexts(ctx context.Context, chatId string) ([]string, error) {
	chatContext, err := c.createChatContext(ctx, chatId)
	if err != nil {
		return nil, err
	}

	return chatContext.getTexts(), nil
}

func (c *Client) Reset(ctx context.Context, chatId string) error {
	chatContext, err := c.createChatContext(ctx, chatId)
	if err != nil {
		return err
	}

	chatContext.reset()
	return c.savechatContext(ctx, chatContext)
}
