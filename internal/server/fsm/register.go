package fsm

func (f *FSM) Register(
	command string,
	handler handler,
	opts... nodeOption,
) {
	n := new(command, handler)
	for _, opt := range opts {
		opt(n)
	}

	f.nodes[command] = n
}

func (f *FSM) Default(
	handler handler,
	opts... nodeOption,
) {
	n := new("default", handler)
	for _, opt := range opts {
		opt(n)
	}

	f.nodes["default"] = n
}
