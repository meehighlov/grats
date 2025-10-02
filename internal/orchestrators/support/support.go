package support

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (o *Orchestrator) SupportHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.SupportHandler(ctx, update)
}

func (o *Orchestrator) WriteHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.WriteHandler(ctx, update)
}

func (o *Orchestrator) CancelHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.CancelHandler(ctx, update)
}

func (o *Orchestrator) SendMessageHandler(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.SendMessageHandler(ctx, update)
}

func (o *Orchestrator) HandleSupportReply(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.HandleSupportReply(ctx, update)
}
