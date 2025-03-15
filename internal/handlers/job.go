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

func notify(ctx context.Context, client *telegram.Client, friends []*db.Friend, chatIdToChat map[string]*db.Chat, logger *slog.Logger) error {
	for _, friend := range friends {
		template := "🔔Сегодня день рождения у %s🥳"
		silent := true
		chat, ok := chatIdToChat[friend.ChatId]
		if ok {
			template = chat.GreetingTemplate
			silent = chat.GetSilent()
		} else {
			logger.Error("Notify job", "error getting chat", "silent mode will be used, default template will be used", "chatid", friend.ChatId)
		}

		msg := fmt.Sprintf(template, friend.Name)

		var sendOpts []telegram.SendMessageOption
		if silent {
			sendOpts = append(sendOpts, telegram.WithDisableNotification())
		}

		logger.Debug("Notify job sending message", "silent", silent)

		_, err := client.SendMessage(ctx, friend.ChatId, msg, sendOpts...)
		if err != nil {
			logger.Error("Notify job", "Notification not sent", err.Error())
			continue
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
			tx.Rollback()
			continue
		}

		if len(friends) == 0 {
			logger.Debug("job", "0 friends found", "continue")
			tx.Rollback()
			continue
		}

		// at most once policy

		isFailed := updateNotifyAt(ctx, tx, friends, logger)
		if isFailed {
			tx.Rollback()
			continue
		}

		chatIdToChat, err := getChatIdToChat(ctx, tx, friends, logger)
		if err != nil {
			logger.Error("Notify job", "error getting chat id to chat", err.Error())
			tx.Rollback()
			continue
		}

		err = commit(ctx, tx, client, logger, cfg)
		if err != nil {
			tx.Rollback()
			continue
		}

		notify(ctx, client, friends, chatIdToChat, logger)
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

func updateNotifyAt(ctx context.Context, tx *sql.Tx, friends []*db.Friend, logger *slog.Logger) bool {
	failed := false
	for _, friend := range friends {
		friend.UpdateNotifyAt()
		err := friend.Save(ctx, tx)
		if err != nil {
			logger.Error("Notify job", "error updating notify date", err.Error(), "chatid", friend.ChatId)
			failed = true
			continue
		}
	}

	return failed
}

func getChatIdToChat(ctx context.Context, tx *sql.Tx, friends []*db.Friend, logger *slog.Logger) (map[string]*db.Chat, error) {
	chatIdToChat := make(map[string]*db.Chat)
	for _, friend := range friends {
		chats, err := (&db.Chat{ChatId: friend.ChatId}).Filter(ctx, tx)
		if err != nil {
			logger.Error("Notify job", "error getting chat", err.Error(), "chatid", friend.ChatId)
			continue
		}
		chatIdToChat[friend.ChatId] = chats[0]
	}
	return chatIdToChat, nil
}
