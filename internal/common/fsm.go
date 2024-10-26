package common

import (
	"context"
	"database/sql"
	"log/slog"
)

func FSM(logger *slog.Logger, handlers map[string]CommandStepHandler) HandlerType {
	return func(ctx context.Context, event Event, tx *sql.Tx) error {
		chatContext := event.GetContext()
		stepTODO := chatContext.GetStepTODO()
		chatContext.SetCommandInProgress(event.GetCommand())

		nextStep := STEPS_DONE

		stepHandler, found := handlers[stepTODO]

		if !found {
			logger.Error("FSM: handler not found, resetting context", "step", stepTODO, "command", event.GetCommand())
			chatContext.Reset()
			return nil
		}

		nextStep, err := stepHandler(ctx, event, tx)
		if err != nil {
			logger.Error("FSM got error from step handler, resseting chat context", "step", stepTODO, "command", event.GetCommand())
			chatContext.Reset()
			return err
		}

		if nextStep == STEPS_DONE {
			chatContext.Reset()
			return nil
		}

		chatContext.SetStepTODO(nextStep)

		return nil
	}
}
