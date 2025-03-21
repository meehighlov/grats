package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

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

	db.MustSetup("grats.db", logger)

	go handlers.BirthdayNotifer(context.Background(), lib.MustSetupLogging("grats_job.log", false, cfg.ENV), cfg)

	updateHandlers := map[string]common.HandlerType{
		// user
		"/start":                              handlers.StartHandler,
		fmt.Sprintf("/start@%s", cfg.BotName): handlers.StartFromGroupHandler,

		"/setup":                              handlers.SetupHandler,
		fmt.Sprintf("/setup@%s", cfg.BotName): handlers.SetupFromGroupHandler,

		"add_to_chat":     handlers.AddToChatHandler,
		"add_enter_bd":    handlers.EnterBirthday,
		"add_save_friend": handlers.SaveFriend,

		// admin TODO

		// support
		"support":               handlers.SupportHandler,
		"write_to_support":      handlers.WriteToSupportHandler,
		"send_to_support":       handlers.SendToSupportHandler,
		"send_support_response": handlers.SendSupportResponseToUserHandler,

		// callback query handlers
		"setup":                       handlers.SetupHandler,
		"list":                        handlers.ListPaginationCallbackQueryHandler,
		"new_list":                    handlers.ListBirthdaysHandler,
		"info":                        handlers.FriendInfoCallbackQueryHandler,
		"delete":                      handlers.DeleteFriendCallbackQueryHandler,
		"confirm_delete":              handlers.ConfirmDeleteFriendCallbackQueryHandler,
		"chat_info":                   handlers.GroupInfoHandler,
		"chat_howto":                  handlers.GroupHowtoHandler,
		"chat_list":                   handlers.GroupHandler,
		"chat_birthdays":              handlers.ListBirthdaysHandler,
		"delete_chat":                 handlers.DeleteChatHandler,
		"confirm_delete_chat":         handlers.ConfirmDeleteChatHandler,
		"edit_greeting_template":      handlers.EditGreetingTemplateHandler,
		"save_greeting_template":      handlers.SaveGreetingTemplateHandler,
		"toggle_silent_notifications": handlers.ToggleSilentNotificationsHandler,
		"edit_name":                   handlers.EditNameHandler,
		"edit_birthday":               handlers.EditBirthdayHandler,
		"save_edit_name":              handlers.SaveEditNameHandler,
		"save_edit_birthday":          handlers.SaveEditBirthdayHandler,
	}

	rootHandler := common.CreateRootHandler(
		logger,
		updateHandlers,
	)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

	var webhookServer *telegram.WebhookServer

	if cfg.UseWebhook {
		logger.Info("starting webhook server...", "addr", cfg.WebhookAddr)
		webhookServer = telegram.StartWebhook(
			cfg.WebhookAddr,
			cfg.BotToken,
			cfg.WebhookSecretToken,
			rootHandler,
			logger,
			cfg.WebhookTlsOn,
			cfg.WebhookTlsCertFile,
			cfg.WebhookTlsKeyFile,
			cfg.WebhookTlsAddr,
		)
	} else {
		logger.Info("starting polling...")
		go telegram.StartPolling(cfg.BotToken, rootHandler)
	}

	<-signalChan
	logger.Info("Sigterm received, stopping...")

	if webhookServer != nil {
		webhookServer.Stop()
	}
}
