package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/telegram"
)

const CHECK_TIMEOUT_SEC = 10

func notify(ctx context.Context, client telegram.ApiCaller, friends []db.Friend, logger *slog.Logger) error {
	msgTemplate := "🔔Сегодня день рождения у %s🥳"
	for _, friend := range friends {
		msg := fmt.Sprintf(msgTemplate, friend.Name)
		_, err := client.SendMessage(ctx, friend.GetChatIdStr(), msg)
		if err != nil {
			logger.Error("Notification not sent:" + err.Error())
		}

		friend.UpdateNotifyAt()
		friend.Save(ctx)
	}

	return nil
}

func run(ctx context.Context, client telegram.ApiCaller, logger *slog.Logger, cfg *config.Config) {
	logger.Info("Starting job for checking birthdays")

	location, err := time.LoadLocation(cfg.Timezone)
	if err != nil {
		panic(err.Error())
	}

	for {
		date := time.Now().In(location).Format("02.01.2006")

		friends, err := (&db.Friend{FilterNotifyAt: date}).Filter(ctx)

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
				logger.Error("panic report error:" + err.Error())
			}
		}
	}()

	run(withCancel, client, logger, cfg)

	return nil
}
