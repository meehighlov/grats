package state

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/meehighlov/grats/internal/repositories/models"
	"github.com/meehighlov/grats/pkg/telegram/fsm"
)

var ErrStateNotFound = errors.New("state not found")

type StoreResult struct {
	value   string
	isFound bool
}

func (s StoreResult) Value() string {
	return s.value
}

func (s StoreResult) IsFound() bool {
	return s.isFound
}

func (r *Repository) GetState(ctx context.Context, key string) (fsm.StoreResult, error) {
	state, err := r.createState(ctx, key)
	isfound := true
	if err == ErrStateNotFound {
		isfound = false
	}
	if err != nil {
		return StoreResult{
			value:   "",
			isFound: isfound,
		}, err
	}
	return StoreResult{
		value:   state.State,
		isFound: isfound,
	}, nil
}

func (r *Repository) SetState(ctx context.Context, key, value string) error {
	state, err := r.createState(ctx, key)
	if err != nil {
		return err
	}
	state.State = value
	return r.saveState(ctx, state)
}

func newState(key string) *models.State {
	return &models.State{
		Key:   key,
		State: "",
	}
}

func (r *Repository) createState(ctx context.Context, key string) (*models.State, error) {
	val, err := r.redis.Redis.Get(ctx, key).Result()

	if err == nil {
		var state models.State
		if err := json.Unmarshal([]byte(val), &state); err == nil {
			return &state, nil
		}
	}

	newState := newState(key)

	jsonState, _ := json.Marshal(newState)
	cmd := r.redis.Redis.Set(ctx, key, jsonState, r.redis.CacheExpiration)
	_, err = cmd.Result()
	if err != nil {
		return &models.State{}, err
	}

	return newState, ErrStateNotFound
}

func (r *Repository) saveState(ctx context.Context, state *models.State) error {
	jsonState, err := json.Marshal(state)
	if err != nil {
		return err
	}

	return r.redis.Redis.Set(ctx, state.Key, jsonState, r.redis.CacheExpiration).Err()
}
