package fsm

func (f *FSM) AddMiddleware(middleware Action) error {
	if middleware == nil {
		return nil
	}

	f.middlewares = append(f.middlewares, middleware)

	return nil
}
