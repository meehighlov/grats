package main

import (
	"context"

	"github.com/meehighlov/grats/internal/auth"
	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/internal/handlers"
	"github.com/meehighlov/grats/internal/lib"
	"github.com/meehighlov/grats/telegram"
)

func main() {
	cfg := config.MustLoad()

	logger := lib.MustSetupLogging("grats.log", true, cfg.ENV)

	// todo use db migrations
	db.MustSetup("grats.db", logger)

	go handlers.BirthdayNotifer(context.Background(), lib.MustSetupLogging("grats_job.log", false, cfg.ENV), cfg)

	updateHandlers := map[string]common.HandlerType {
		// user
		"/start": auth.Auth(logger, handlers.StartHandler),
		"/help": auth.Auth(logger, handlers.HelpHandler),
		"/list": auth.Auth(logger, handlers.ListBirthdaysHandler),
		"/add": auth.Auth(logger, common.FSM(logger, handlers.AddBirthdayChatHandler())),

		// admin
		"/access_list": auth.Admin(logger, handlers.AccessListHandler),
		"/access_grant": auth.Admin(logger, common.FSM(logger, handlers.GrantAccessChatHandler())),
		"/access_revoke": auth.Admin(logger, common.FSM(logger, handlers.RevokeAccessChatHandler())),

		// callback query handlers
		"list": handlers.ListBirthdaysCallbackQueryHandler,
		"info": handlers.FriendInfoCallbackQueryHandler,
		"delete": handlers.DeleteFriendCallbackQueryHandler,
	}

	rootHandler := common.CreateRootHandler(
		logger,
		common.NewChatCache(),
		updateHandlers,
	)

	logger.Info("starting polling...")
	telegram.StartPolling(cfg.BotToken, rootHandler)
}
