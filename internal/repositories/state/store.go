package state

import (
	"context"
	"encoding/json"

	"github.com/meehighlov/grats/internal/repositories/entities"
)

func (r *Repository) GetState(ctx context.Context, key string) (string, error) {
	state, err := r.createState(ctx, key)
	if err != nil {
		return "", err
	}
	return state.State, nil
}

func (r *Repository) SetState(ctx context.Context, key, value string) error {
	state, err := r.createState(ctx, key)
	if err != nil {
		return err
	}
	state.State = value
	return r.saveState(ctx, state)
}

func newState(key string) *entities.State {
	return &entities.State{
		Key:    key,
		State:  "",
	}
}

func (r *Repository) createState(ctx context.Context, key string) (*entities.State, error) {
	val, err := r.redis.Redis.Get(ctx, key).Result()

	if err == nil {
		var state entities.State
		if err := json.Unmarshal([]byte(val), &state); err == nil {
			return &state, nil
		}
	}

	newState := newState(key)

	jsonState, _ := json.Marshal(newState)
	cmd := r.redis.Redis.Set(ctx, key, jsonState, r.redis.CacheExpiration)
	_, err = cmd.Result()
	if err != nil {
		return &entities.State{}, err
	}

	return newState, nil
}

func (r *Repository) saveState(ctx context.Context, state *entities.State) error {
	jsonState, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return r.redis.Redis.Set(ctx, state.Key, jsonState, r.redis.CacheExpiration).Err()
}
