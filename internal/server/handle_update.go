package server

import (
	"context"
	"runtime/debug"
	"strings"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
)

func (s *Server) HandleUpdate(ctx context.Context, update *telegram.Update) error {
	if update.IsCallback() {
		s.clients.Telegram.AnswerCallbackQuery(ctx, update.CallbackQuery.Id)
	}

	if err := s.clients.Cache.CreateChatContext(ctx, update.GetChatIdStr()); err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			s.logger.Error(
				"Root handler",
				"recovered from panic, error", r,
				"stack", string(debug.Stack()),
				"update", update,
			)
			s.clients.Cache.Reset(ctx, update.GetChatIdStr())

			chatId := update.GetChatIdStr()
			if chatId != "" {
				s.clients.Telegram.Reply(ctx, s.constants.ERROR_MESSAGE, update)
				return
			}

			s.logger.Error(
				"Root handler",
				"recover from panic", "chatId was empty",
				"update", update,
			)
		}
	}()

	// Handling support replies
	if update.GetMessage() != nil &&
		update.GetMessage().GetChatIdStr() == s.cfg.SupportChatId &&
		update.GetMessage().IsReply() {
		if err := s.orchestrators.Support.HandleSupportReply(ctx, update); err != nil {
			s.logger.Error("Failed to handle support reply", "error", err)
		}
		return nil
	}

	command_ := update.GetMessage().GetCommand()
	command := ""

	if command_ != "" {
		command = command_
		s.clients.Cache.Reset(ctx, update.GetChatIdStr())
	} else {
		if update.CallbackQuery.Id != "" {
			params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

			s.logger.Info("CallbackQueryHandler", "command", params.Command, "chat id", update.GetChatIdStr())
			command = params.Command
		} else {
			command_, err := s.clients.Cache.GetNextHandler(ctx, update.GetChatIdStr())
			if err != nil {
				return err
			}
			if command_ != "" {
				command = command_
			}
		}
	}

	s.logger.Info("Root handler", "handling", command, "with update", update)

	err := s.handle(ctx, update, command)
	if err != nil {
		s.clients.Cache.Reset(ctx, update.GetChatIdStr())
		s.logger.Error("Root handler", "error", err.Error(), "chat id", update.GetChatIdStr(), "update id", update.UpdateId)
	} else {
		s.logger.Info("Root handler", "success", command, "chat id", update.GetChatIdStr(), "update id", update.UpdateId)
	}

	return nil
}

func (s *Server) handle(ctx context.Context, update *telegram.Update, command string) error {
	switch command {
	// user
	case s.constants.CMD_START:
		return s.orchestrators.User.Start(ctx, update)

	// wish handlers
	case s.constants.CMD_ADD_TO_WISH:
		return s.orchestrators.Wish.AddWishHandler(ctx, update)
	case s.constants.CMD_ADD_SAVE_WISH:
		return s.orchestrators.Wish.SaveWish(ctx, update)
	case s.constants.CMD_WISH_LIST, s.constants.CMD_WISHLIST:
		return s.orchestrators.Wish.List(ctx, update)
	case s.constants.CMD_WISH_INFO:
		return s.orchestrators.Wish.WishInfoHandler(ctx, update)
	case s.constants.CMD_DELETE_WISH:
		return s.orchestrators.Wish.DeleteWishCallbackQueryHandler(ctx, update)
	case s.constants.CMD_CONFIRM_DELETE_WISH:
		return s.orchestrators.Wish.ConfirmDeleteWishCallbackQueryHandler(ctx, update)
	case s.constants.CMD_EDIT_PRICE:
		return s.orchestrators.Wish.EditPriceHandler(ctx, update)
	case s.constants.CMD_EDIT_LINK:
		return s.orchestrators.Wish.EditLinkHandler(ctx, update)
	case s.constants.CMD_EDIT_PRICE_SAVE:
		return s.orchestrators.Wish.SaveEditPriceHandler(ctx, update)
	case s.constants.CMD_EDIT_LINK_SAVE:
		return s.orchestrators.Wish.SaveEditLinkHandler(ctx, update)
	case s.constants.CMD_EDIT_WISH_NAME:
		return s.orchestrators.Wish.EditWishNameHandler(ctx, update)
	case s.constants.CMD_EDIT_WISH_NAME_SAVE:
		return s.orchestrators.Wish.SaveEditWishNameHandler(ctx, update)
	case s.constants.CMD_SHARE_WISH_LIST:
		return s.orchestrators.Wish.ShareWishListHandler(ctx, update)
	case s.constants.CMD_SHOW_SWL:
		return s.orchestrators.Wish.ShowSharedWishlistHandler(ctx, update)
	case s.constants.CMD_SHOW_SWI:
		return s.orchestrators.Wish.WishInfoHandler(ctx, update)
	case s.constants.CMD_TOGGLE_WISH_LOCK:
		return s.orchestrators.Wish.ToggleWishLockHandler(ctx, update)

	// support handlers
	case s.constants.CMD_SUPPORT:
		return s.orchestrators.Support.SupportHandler(ctx, update)
	case s.constants.CMD_SUPPORT_WRITE:
		return s.orchestrators.Support.WriteHandler(ctx, update)
	case s.constants.CMD_SUPPORT_CANCEL:
		return s.orchestrators.Support.CancelHandler(ctx, update)
	case s.constants.CMD_SUPPORT_SEND:
		return s.orchestrators.Support.SendMessageHandler(ctx, update)

	// callback query handlers
	case s.constants.CMD_LIST:
		return s.orchestrators.Wish.List(ctx, update)
	default:
		if strings.HasPrefix(command, s.constants.CMD_START) {
			idForCommand := strings.TrimSpace(strings.TrimPrefix(command, s.constants.CMD_START))

			if idForCommand != "" && strings.HasPrefix(idForCommand, s.constants.SHARED_LIST_ID_PREFIX) {
				return s.orchestrators.Wish.ShowSharedWishlistHandler(ctx, update)
			}
		}

		s.logger.Debug(
			"default case",
			"not recognized command", command,
		)
	}

	return nil
}
