package fsm

import "context"

type StoreResult interface {
	Value() string
	IsFound() bool
}

type StateStore interface {
	GetState(context.Context, string) (StoreResult, error)
	SetState(context.Context, string, string) error
}
