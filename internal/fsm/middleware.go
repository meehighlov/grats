package fsm

import "github.com/meehighlov/grats/internal/fsm/handler"

func (f *FSM) AddMiddleware(handler handler.HandlerType) error {
	if handler == nil {
		return nil
	}

	f.middlewares = append(f.middlewares, handler)

	return nil
}
