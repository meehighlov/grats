package telegram

import (
	"context"
	"errors"

	"github.com/meehighlov/grats/pkg/telegram/fsm"
)

type Handler func(context.Context, *Scope) error
type Condition func(context.Context, *Scope) (bool, error)

func (b *Bot) Serve() error {
	return b.server.Serve()
}

func (b *Bot) AddHandler(
	handler Handler,
	condition Condition,
	options ...fsm.StateOption,
) *fsm.State {
	return b.fsm.Activate(
		wrapHandler(handler),
		wrapCondition(condition),
		options...,
	)
}

func (b *Bot) Reset(
	condition Condition,
	options ...fsm.StateOption,
) *fsm.State {
	return b.fsm.Reset(wrapCondition(condition), options...)
}

func (b *Bot) AddMiddleware(middleware Handler) error {
	return b.fsm.AddMiddleware(wrapHandler(middleware))
}

func BeforeAction(beforeAction Handler) fsm.StateOption {
	return func(s *fsm.State) error {
		s.SetBeforeAction(wrapHandler(beforeAction))
		return nil
	}
}

func AcceptFrom(fromState *fsm.State) fsm.StateOption {
	return func(s *fsm.State) error {
		fromState.AddTransition(s)
		return nil
	}
}

func wrapHandler(handler Handler) fsm.Action {
	return func(ctx context.Context, data fsm.ActionData) error {
		scope, err := castScope(data)
		if err != nil {
			return err
		}
		return handler(ctx, scope)
	}
}

func wrapCondition(condition Condition) fsm.Condition {
	return func(ctx context.Context, data fsm.ActionData) (bool, error) {
		scope, err := castScope(data)
		if err != nil {
			return false, err
		}
		return condition(ctx, scope)
	}
}

func castScope(data fsm.ActionData) (*Scope, error) {
	scope, ok := data.(*Scope)
	if !ok {
		return nil, errors.New("data is not a telegram scope")
	}
	return scope, nil
}
