package fsm

import (
	"context"
	"errors"
	"fmt"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/fsm/state"
)

func (f *FSM) Handle(ctx context.Context, update *telegram.Update) error {
	defer func() error {
		r := recover()
		if r != nil {
			critical := fmt.Errorf("recover from panic: %v", r)
			err := f.stateStore.SetState(ctx, update.GetChatIdStr(), state.READY)
			return errors.Join(critical, err)
		}
		return nil
	}()

	for _, middleware := range f.middlewares {
		if err := middleware(ctx, update); err != nil {
			return err
		}
	}

	status, err := f.stateStore.GetState(ctx, update.GetChatIdStr())
	if err != nil {
		return err
	}

	var s *state.State
	for _, state := range f.states {
		ok, err := state.Condition().Check(ctx, update)
		if err != nil {
			return err
		}

		if ok && state.IsActivationAllowed(status) {
			s = state
			break
		}
	}

	if s == nil {
		return fmt.Errorf("not found handler for update: %v", update)
	}

	err = s.Activate(ctx, update)

	cerr := f.stateStore.SetState(
		ctx,
		update.GetChatIdStr(),
		s.Next(err),
	)

	return errors.Join(err, cerr)
}
