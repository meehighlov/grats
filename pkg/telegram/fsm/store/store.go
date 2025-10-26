package store

import "context"

// TODO make more detailed interface
// return (Result, error)
// Result {state, isfound}
// to check if !isfound -> use ready state as deafult
type StateStore interface {
	GetState(ctx context.Context, key string) (string, error)
	SetState(ctx context.Context, key, value string) error
}
