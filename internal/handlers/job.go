package handlers

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"time"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/telegram"
)

const CHECK_TIMEOUT_SEC = 10

func notify(ctx context.Context, client *telegram.Client, friends []db.Friend, logger *slog.Logger) error {
	msgTemplate := "🔔Сегодня день рождения у %s (%s)🥳"
	for _, friend := range friends {
		msg := fmt.Sprintf(msgTemplate, friend.Name, friend.City)
		_, err := client.SendMessage(ctx, friend.ChatId, msg)
		if err != nil {
			logger.Error("Notify job", "Notification not sent", err.Error())
		}

		tx, err := db.GetDBConnection().BeginTx(ctx, nil)
		if err != nil {
			logger.Error("Notify job", "getting transaction error", err.Error())
			continue
		}

		friend.UpdateNotifyAt()
		err = friend.Save(ctx, tx)
		if err != nil {
			logger.Error("Notify job", "error updating notify date", err.Error())
			tx.Rollback()
		} else {
			tx.Commit()
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
			logger.Debug("Notify job", "error filtering friends", err.Error())
			tx.Rollback()
		} else {
			tx.Commit()
		}

		if err != nil {
			logger.Error("Error getting birthdays: " + err.Error())
		} else {
			notify(ctx, client, friends, logger)
		}

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
