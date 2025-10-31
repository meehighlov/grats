package telegram

import (
	"context"

	"github.com/meehighlov/grats/pkg/telegram/builders"
	"github.com/meehighlov/grats/pkg/telegram/client"
	"github.com/meehighlov/grats/pkg/telegram/fsm"
	"github.com/meehighlov/grats/pkg/telegram/models"
)

type Scoper struct {
	fsm      *fsm.FSM
	client   *client.Client
	builders *builders.Builders
}

func NewScoper(fsm *fsm.FSM, client *client.Client, builders *builders.Builders) *Scoper {
	return &Scoper{
		fsm:      fsm,
		client:   client,
		builders: builders,
	}
}

func (w *Scoper) Handle(ctx context.Context, update *models.Update) error {
	scope := NewScope(w.client, w.builders, update)
	return w.fsm.Handle(ctx, scope)
}
