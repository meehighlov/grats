package main

import (
	"context"

	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/internal/auth"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/handlers"
	"github.com/meehighlov/grats/internal/lib"
	"github.com/meehighlov/grats/telegram"
)

func main() {
	cfg := config.MustLoad()

	logger := lib.MustSetupLogging("grats.log", true, cfg.ENV)

	// todo use db migrations
	db.MustSetup("grats.db", logger)

	bot := telegram.NewBot(cfg.BotToken)

	go handlers.BirthdayNotifer(
		context.Background(),
		cfg.BotToken,
		lib.MustSetupLogging("job.log", false, cfg.ENV),
	)

	bot.RegisterCommandHandler("/start", auth.Auth(handlers.StartHandler))
	bot.RegisterCommandHandler("/help", auth.Auth(handlers.HelpHandler))
	bot.RegisterCommandHandler("/list", auth.Auth(handlers.ListBirthdaysHandler))
	bot.RegisterCommandHandler("/add", auth.Auth(telegram.FSM(handlers.AddBirthdayChatHandler())))
	bot.RegisterCommandHandler("/access_list", auth.Admin(handlers.AccessListHandler))
	bot.RegisterCommandHandler("/access_grant", auth.Admin(telegram.FSM(handlers.GrantAccessChatHandler())))
	bot.RegisterCommandHandler("/access_revoke", auth.Admin(telegram.FSM(handlers.RevokeAccessChatHandler())))

	bot.StartPolling()
	logger.Info("Polling started.")
}
