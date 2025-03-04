package main

import (
	"context"

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
		"/start": auth.Auth(logger, handlers.StartHandler),
		"/list":  auth.Auth(logger, handlers.ListBirthdaysHandler),

		"/add":            auth.Auth(logger, handlers.AddToPrivateListHandler),
		"add_to_chat":     handlers.AddToChatHandler,
		"add_enter_bd":    handlers.EnterBirthday,
		"add_save_friend": handlers.SaveFriend,

		"/chats": auth.Auth(logger, handlers.GroupHandler),

		// admin
		"/admin":                  auth.Admin(logger, admin.AdminCommandListHandler),
		"/access_list":            auth.Admin(logger, admin.AccessListHandler),
		"/access_grant":           auth.Admin(logger, admin.GrantAccess),
		"access_save_tg_username": admin.SaveAccess,
		"/access_revoke":          auth.Admin(logger, admin.RevokeAccess),
		"access_update":           admin.UpdateAccessInfo,

		// callback query handlers
		"list":                   handlers.ListPaginationCallbackQueryHandler,
		"info":                   handlers.FriendInfoCallbackQueryHandler,
		"delete":                 handlers.DeleteFriendCallbackQueryHandler,
		"chat_info":              handlers.GroupInfoHandler,
		"chat_howto":             handlers.GroupHowtoHandler,
		"chat_list":              handlers.GroupHandler,
		"chat_birthdays":         handlers.ListBirthdaysHandler,
		"chat_delete":            handlers.GroupChatRegisterHandler,
		"edit_greeting_template": handlers.EditGreetingTemplateHandler,
		"save_greeting_template": handlers.SaveGreetingTemplateHandler,

		// group chat handler
		"chat_register": handlers.GroupChatRegisterHandler,
	}

	rootHandler := common.CreateRootHandler(
		logger,
		updateHandlers,
	)

	logger.Info("starting polling...")
	telegram.StartPolling(cfg.BotToken, rootHandler)
}
