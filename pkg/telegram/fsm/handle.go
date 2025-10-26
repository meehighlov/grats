package fsm

import (
	"context"
	"errors"
	"fmt"

	"github.com/meehighlov/grats/pkg/telegram/fsm/state"
	"github.com/meehighlov/grats/pkg/telegram/models"
)

func (f *FSM) Handle(ctx context.Context, update *models.Update) error {
	defer func() error {
		r := recover()
		if r != nil {
			critical := fmt.Errorf("recover from panic: %v", r)
			key := update.GetChatIdStr() + ":state"
			err := f.stateStore.SetState(ctx, key, state.READY.String())
			return errors.Join(critical, err)
		}
		return nil
	}()

	for _, middleware := range f.middlewares {
		if err := middleware(ctx, update); err != nil {
			return err
		}
	}

	key := update.GetChatIdStr() + ":state"

	currentStateId, err := f.stateStore.GetState(ctx, key)
	if err != nil {
		return err
	}

	var s *state.State
	for _, state := range f.states[currentStateId].GetTransitions() {
		ok, err := state.Condition().Check(ctx, update)
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

	err = s.Activate(ctx, update)

	cerr := f.stateStore.SetState(
		ctx,
		key,
		s.Done(err, currentStateId),
	)

	return errors.Join(err, cerr)
}

func (f *FSM) reset(ctx context.Context, update *models.Update) error {
	key := update.GetChatIdStr() + ":state"
	return f.stateStore.SetState(ctx, key, state.READY.String())
}
