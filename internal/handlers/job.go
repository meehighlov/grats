package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
)

const CHECK_TIMEOUT_SEC = 10

func notify(ctx context.Context, client *telegram.Client, friends []db.Friend, logger *slog.Logger, tx *sql.Tx) error {
	for _, friend := range friends {
		template := "🔔Сегодня день рождения у %s🥳"
		chat := db.Chat{}
		chat.BaseFields.ID = friend.ChatId
		chats, err := chat.Filter(ctx, tx)
		if err != nil {
			logger.Error("Notify job", "error getting chat, default template will be used", err.Error())
		}

		if len(chats) > 0 && chats[0].GreetingTemplate != "" {
			template = chats[0].GreetingTemplate
		}

		msg := fmt.Sprintf(template, friend.Name)
		_, err = client.SendMessage(ctx, chats[0].TGChatId, msg)
		if err != nil {
			logger.Error("Notify job", "Notification not sent", err.Error())
			continue
		}

		friend.UpdateNotifyAt()
		err = friend.Save(ctx, tx)
		if err != nil {
			logger.Error("Notify job", "error updating notify date", err.Error(), "chatid", friend.ChatId)
		}
	}

	return nil
}

func run(ctx context.Context, client *telegram.Client, logger *slog.Logger, cfg *config.Config) {
	logger.Info("Starting job for checking birthdays")

	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		log.Fatal(err.Error())
	}

	for {
		date := time.Now().In(location).Format("02.01.2006")

		tx, err := db.GetDBConnection().BeginTx(ctx, nil)
		if err != nil {
			logger.Error("Notify job", "getting transaction error", err.Error())
			continue
		}

		friends, err := (&db.Friend{FilterNotifyAt: date}).Filter(ctx, tx)
		logger.Info("Notify job", "found rows", len(friends))
		if err != nil {
			logger.Debug("Notify job", "Error getting birthdays", err.Error())
			continue
		}

		notify(ctx, client, friends, logger, tx)

		tx.Commit()

		time.Sleep(CHECK_TIMEOUT_SEC * time.Second)
	}
}

func BirthdayNotifer(
	ctx context.Context,
	logger *slog.Logger,
	cfg *config.Config,
) error {
	withCancel, cancel := context.WithCancel(ctx)
	defer cancel()

	client := telegram.NewClient(cfg.BotToken, logger)

	defer func() {
		if r := recover(); r != nil {
			errMsg := fmt.Sprintf("Паника в фоновой задаче проверки дней рождения\n %s", r)

			logger.Error(errMsg)

			reportChatId := cfg.ReportChatId
			_, err := client.SendMessage(context.Background(), reportChatId, errMsg)
			if err != nil {
				logger.Error("report fatal error:" + err.Error())
			}
		}
	}()

	run(withCancel, client, logger, cfg)

	return nil
}
