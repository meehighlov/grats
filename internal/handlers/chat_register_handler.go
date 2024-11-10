package handlers

import (
	"context"
	"database/sql"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
)

func GroupChatRegisterHandler(ctx context.Context, event *common.Event, tx *sql.Tx) error {
	chat := db.Chat{
		BaseFields:   db.NewBaseFields(),
		ChatType:     "group",
		BotInvitedBy: strconv.Itoa(event.GetMessage().From.Id),
		ChatId:       event.GetMessage().GetChatIdStr(),
	}

	message := event.GetMessage()

	if message.LeftChatMember.Username == config.Cfg().BotName {
		// todo check bot name by GetMe tg method
		err := chat.Delete(ctx, tx)
		if err != nil {
			return err
		}
		err = (&db.Friend{ChatId: strconv.Itoa(message.Chat.Id)}).Delete(ctx, tx)
		if err != nil {
			return err
		}
		return nil
	}

	err := chat.Save(ctx, tx)
	if err != nil {
		return err
	}

	return nil
}
