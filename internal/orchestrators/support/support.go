package support

import (
	"context"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (o *Orchestrator) Support(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.Support(ctx, update)
}

func (o *Orchestrator) SupportWrite(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.SupportWrite(ctx, update)
}

func (o *Orchestrator) CancelSupportCall(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.CancelSupportCall(ctx, update)
}

func (o *Orchestrator) SendSupportMessage(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.SendSupportMessage(ctx, update)
}

func (o *Orchestrator) ProcessSupportReply(ctx context.Context, update *telegram.Update) error {
	return o.services.Support.ProcessSupportReply(ctx, update)
}
