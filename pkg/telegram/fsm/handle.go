package fsm

import (
	"context"
	"errors"
	"fmt"
)

func (f *FSM) Handle(ctx context.Context, data ActionData) error {
	defer func() error {
		r := recover()
		if r != nil {
			critical := fmt.Errorf("recover from panic: %v", r)
			key := data.UserID() + ":state"
			err := f.stateStore.SetState(ctx, key, READY.String())
			return errors.Join(critical, err)
		}
		return nil
	}()

	for _, middleware := range f.middlewares {
		if err := middleware(ctx, data); err != nil {
			return err
		}
	}

	key := data.UserID() + ":state"

	currentStateId, err := f.stateStore.GetState(ctx, key)
	if err != nil {
		return err
	}

	var s *State
	for _, state := range f.states[currentStateId].GetTransitions() {
		ok, err := state.Condition(ctx, data)
		if err != nil {
			return err
		}

		if ok {
			s = state
			break
		}
	}

	if s == nil {
		return fmt.Errorf("not found transition for state %s", currentStateId)
	}

	err = s.Activate(ctx, data)

	cerr := f.stateStore.SetState(
		ctx,
		key,
		s.Done(err, currentStateId),
	)

	return errors.Join(err, cerr)
}

func (f *FSM) reset(ctx context.Context, data ActionData) error {
	userId := data.UserID()
	key := userId + ":state"
	return f.stateStore.SetState(ctx, key, READY.String())
}
