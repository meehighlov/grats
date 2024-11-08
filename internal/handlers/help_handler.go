package handlers

import (
	"context"
	"database/sql"
	"strings"

	"github.com/meehighlov/grats/internal/common"
)

func AdminCommandListHandler(ctx context.Context, event common.Event, tx *sql.Tx) error {
	commands := []string{
		"/access_list - список пользователей с доступом😏",
		"/access_grant - предоставить доступ🙈",
		"/access_revoke - отозвать доступ🤝",
	}

	msg := strings.Join(commands, "\n")

	event.Reply(ctx, msg)

	return nil
}
