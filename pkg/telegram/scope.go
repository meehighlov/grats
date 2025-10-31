package telegram

import (
	"github.com/meehighlov/grats/pkg/telegram/builders"
	"github.com/meehighlov/grats/pkg/telegram/client"
	"github.com/meehighlov/grats/pkg/telegram/models"
)

type Scope struct {
	client   *client.Client
	builders *builders.Builders
	update   *models.Update
}

func NewScope(client *client.Client, builders *builders.Builders, update *models.Update) *Scope {
	return &Scope{
		client:   client,
		builders: builders,
		update:   update,
	}
}
