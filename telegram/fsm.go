package telegram

import (
	"fmt"
	"log"
)

func FSM(handlers map[int]CommandStepHandler) CommandHandler {
	return func(event Event) error {
		ctx := event.GetContext()
		stepTODO := ctx.getStepTODO()
		ctx.setCommandInProgress(event.getCommand())

		nextStep := STEPS_DONE

		stepHandler, found := handlers[stepTODO]

		if !found {
			logMsg := fmt.Sprintf("Step %d not supported for %s, resetting context", stepTODO, event.getCommand())
			log.Println(logMsg)
			ctx.reset()
			return nil
		}

		nextStep, _ = stepHandler(event)

		if nextStep == STEPS_DONE {
			ctx.reset()
			return nil
		}

		ctx.setStepTODO(nextStep)

		return nil
	}
}
