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
const COMMIT_EXECUTION_TIMEOUT_SEC = 10

func notify(ctx context.Context, client *telegram.Client, friends []*db.Friend, logger *slog.Logger, tx *sql.Tx, cfg *config.Config) error {
	for _, friend := range friends {
		template := "🔔Сегодня день рождения у %s🥳"
		chats, err := (&db.Chat{ChatId: friend.ChatId}).Filter(ctx, tx)
		if err != nil {
			logger.Error("Notify job", "error getting chat, default template will be used", err.Error())
		}

		if len(chats) == 0 {
			logger.Error("Notify job", "error getting chat, skipping notification", err.Error())
			continue
		}

		if chats[0].GreetingTemplate != "" {
			template = chats[0].GreetingTemplate
		}

		msg := fmt.Sprintf(template, friend.Name)

		var sendOpts []telegram.SendMessageOption
		if chats[0].IsAlreadySilent() {
			sendOpts = append(sendOpts, telegram.WithDisableNotification())
		}

		_, err = client.SendMessage(ctx, friend.ChatId, msg, sendOpts...)
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
		time.Sleep(CHECK_TIMEOUT_SEC * time.Second)

		date := time.Now().In(location).Format("02.01.2006")

		tx, err := db.GetDBConnection().BeginTx(ctx, nil)
		if err != nil {
			logger.Error("Notify job", "getting transaction error", err.Error())
			continue
		}

		friends, err := (&db.Friend{FilterNotifyAt: date}).Filter(ctx, tx)
		if err != nil {
			logger.Error("Notify job", "error getting notification list", err.Error())
			continue
		}

		// renew notify_at for friends
		// and then commit
		// and then notify

		notify(ctx, client, friends, logger, tx, cfg)

		commitContext, cancel := context.WithTimeout(ctx, COMMIT_EXECUTION_TIMEOUT_SEC*time.Second)
		defer cancel()
		err = commit(commitContext, tx, client, logger, cfg)
		if err != nil {
			break
		}
		cancel()
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

func commit(
	ctx context.Context,
	tx *sql.Tx,
	client *telegram.Client,
	logger *slog.Logger,
	cfg *config.Config,
) error {
	err := tx.Commit()
	if err == nil {
		return nil
	}

	errMsg := fmt.Sprintf(
		"Notify job: Не удалось установить комит для notify_at: %s",
		err.Error(),
	)
	logger.Error(
		"Notify job", "error committing transaction",
		err.Error(),
	)
	_, err = client.SendMessage(
		ctx,
		cfg.ReportChatId,
		errMsg,
	)
	if err != nil {
		logger.Error("Notify job", "error sending report message", err.Error())
	}

	return err
}
