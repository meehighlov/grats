package common

import (
	"log/slog"
)

func FSM(logger *slog.Logger, handlers map[string]CommandStepHandler) HandlerType {
	return func(event Event) error {
		ctx := event.GetContext()
		stepTODO := ctx.GetStepTODO()
		ctx.SetCommandInProgress(event.GetCommand())

		nextStep := STEPS_DONE

		stepHandler, found := handlers[stepTODO]

		if !found {
			logger.Error("FSM: handler not found, resetting context", "step", stepTODO, "command", event.GetCommand())
			ctx.Reset()
			return nil
		}

		nextStep, _ = stepHandler(event)

		if nextStep == STEPS_DONE {
			ctx.Reset()
			return nil
		}

		ctx.SetStepTODO(nextStep)

		return nil
	}
}
