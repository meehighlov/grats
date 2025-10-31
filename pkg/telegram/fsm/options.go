package fsm

func BeforeAction(beforeAction Action) StateOption {
	return func(s *State) error {
		if beforeAction != nil {
			s.SetBeforeAction(beforeAction)
		}
		return nil
	}
}

func AcceptFrom(fromState *State) StateOption {
	return func(s *State) error {
		fromState.AddTransition(s)
		return nil
	}
}
