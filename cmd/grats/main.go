package main

import (
	"context"
	"fmt"

	"github.com/meehighlov/grats/internal/auth"
	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/internal/handlers"
	"github.com/meehighlov/grats/internal/handlers/admin"
	"github.com/meehighlov/grats/internal/lib"
	"github.com/meehighlov/grats/telegram"
)

func main() {
	cfg := config.MustLoad()

	logger := lib.MustSetupLogging("grats.log", true, cfg.ENV)

	db.MustSetup("grats.db", logger)

	go handlers.BirthdayNotifer(context.Background(), lib.MustSetupLogging("grats_job.log", false, cfg.ENV), cfg)

	updateHandlers := map[string]common.HandlerType{
		// user
		"/start":                              auth.Auth(logger, handlers.StartHandler),
		fmt.Sprintf("/start@%s", cfg.BotName): auth.Auth(logger, handlers.StartFromGroupHandler),

		"/setup":                              auth.Auth(logger, handlers.SetupHandler),
		fmt.Sprintf("/setup@%s", cfg.BotName): auth.Auth(logger, handlers.SetupFromGroupHandler),

		"add_to_chat":     handlers.AddToChatHandler,
		"add_enter_bd":    handlers.EnterBirthday,
		"add_save_friend": handlers.SaveFriend,

		// admin
		"/admin":                  auth.Admin(logger, admin.AdminCommandListHandler),
		"/access_list":            auth.Admin(logger, admin.AccessListHandler),
		"/access_grant":           auth.Admin(logger, admin.GrantAccess),
		"access_save_tg_username": admin.SaveAccess,
		"/access_revoke":          auth.Admin(logger, admin.RevokeAccess),
		"access_update":           admin.UpdateAccessInfo,

		// callback query handlers
		"setup":                  handlers.SetupHandler,
		"list":                   handlers.ListPaginationCallbackQueryHandler,
		"new_list":               handlers.ListBirthdaysHandler,
		"info":                   handlers.FriendInfoCallbackQueryHandler,
		"delete":                 handlers.DeleteFriendCallbackQueryHandler,
		"confirm_delete":         handlers.ConfirmDeleteFriendCallbackQueryHandler,
		"chat_info":              handlers.GroupInfoHandler,
		"chat_howto":             handlers.GroupHowtoHandler,
		"chat_list":              handlers.GroupHandler,
		"chat_birthdays":         handlers.ListBirthdaysHandler,
		"delete_chat":            handlers.DeleteChatHandler,
		"confirm_delete_chat":    handlers.ConfirmDeleteChatHandler,
		"edit_greeting_template": handlers.EditGreetingTemplateHandler,
		"save_greeting_template": handlers.SaveGreetingTemplateHandler,
	}

	rootHandler := common.CreateRootHandler(
		logger,
		updateHandlers,
	)

	logger.Info("starting polling...")
	telegram.StartPolling(cfg.BotToken, rootHandler)
}
