package fsm

import "github.com/meehighlov/grats/pkg/telegram/fsm/action"

func (f *FSM) AddMiddleware(action action.Action) error {
	if action == nil {
		return nil
	}

	f.middlewares = append(f.middlewares, action)

	return nil
}
