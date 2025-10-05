package wish

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	inlinekeyboard "github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/entities"
	"github.com/meehighlov/grats/internal/repositories/wish"
)

func (s *Service) ShareWishListHandler(ctx context.Context, update *telegram.Update) error {
	params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)

	wishListId := params.ID

	wishlist, err := s.repositories.WishList.Filter(ctx, &entities.WishList{BaseFields: entities.BaseFields{ID: wishListId}})
	if err != nil {
		return err
	}

	if len(wishlist) == 0 {
		return fmt.Errorf("wishlist not found")
	}

	botName := s.cfg.BotName
	shareLink := fmt.Sprintf(s.constants.SHARE_WISHLIST_LINK_TEMPLATE, botName, wishlist[0].ID)

	keyboard := s.builders.KeyboardBuilder.NewKeyboard()
	keyboard.AppendAsLine(
		keyboard.NewShareLinkButton(s.constants.BTN_SHARE, shareLink, s.constants.MY_WISHLIST_SHARE_TITLE),
		keyboard.NewCopyButton(s.constants.BTN_COPY_LINK, shareLink),
	)
	keyboard.AppendAsLine(
		keyboard.NewButton(s.constants.BTN_BACK_TO_WISHLIST, s.builders.CallbackDataBuilder.Build(wishlist[0].ID, s.constants.CMD_LIST, s.constants.LIST_DEFAULT_OFFSET).String()),
	)

	shareMessage := s.constants.SHARE_WISHLIST_MESSAGE

	if _, err := s.clients.Telegram.Edit(
		ctx,
		shareMessage,
		update,
		telegram.WithReplyMurkup(keyboard.Murkup()),
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) ShowSharedWishlistHandler(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()

	listPrefix := s.constants.CMD_START + " " + s.constants.SHARED_LIST_ID_PREFIX

	wishlistId := strings.TrimPrefix(message.Text, listPrefix)
	offset := s.constants.LIST_START_OFFSET
	if update.IsCallback() {
		params := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)
		wishlistId = params.ID
		offset, _ = strconv.Atoi(params.Offset)
		if offset == 0 {
			offset = s.constants.LIST_START_OFFSET
		}
	}

	wishes, err := s.repositories.Wish.List(ctx, &wish.ListFilter{WishListID: wishlistId, Limit: s.cfg.ListLimitLen, Offset: offset})
	if err != nil {
		s.logger.Error(
			"StartHandler - Shared wishlist",
			"error", "error getting wishes",
			"details", err.Error(),
			"wl_id", wishlistId,
			"chatId", message.GetChatIdStr(),
		)
		s.clients.Telegram.Reply(ctx, s.constants.FAILED_TO_LOAD_WISHES, update)
		return err
	}

	if len(wishes) == 0 {
		// TODO: update message when user list is empty
		return nil
	}

	// user used his own link
	// just send greeting
	if s.cfg.IsProd() {
		if wishes[0].UserId == strconv.Itoa(message.From.Id) {
			s.clients.Telegram.Reply(ctx, s.constants.HELLO_AGAIN, update)
			return nil
		}
	}

	userInfo, _ := s.clients.Telegram.GetChatMember(ctx, wishes[0].UserId)
	name := userInfo.Result.User.Username
	if name == "" {
		name = userInfo.Result.User.FirstName
	}
	header := fmt.Sprintf(s.constants.WISHLIST_HEADER_TEMPLATE, "@"+name)

	count, err := s.repositories.Wish.Count(ctx, &wish.CountFilter{WishListID: wishlistId})
	if err != nil {
		return err
	}

	keyboard := s.buildSharedWishlistMarkup(wishes, int(count), offset, wishlistId)

	if update.IsCallback() {
		if _, err := s.clients.Telegram.Edit(
			ctx,
			header,
			update,
			telegram.WithReplyMurkup(keyboard.Murkup()),
		); err != nil {
			return err
		}

		return nil
	}

	if _, err := s.clients.Telegram.Reply(
		ctx,
		header,
		update,
		telegram.WithReplyMurkup(keyboard.Murkup()),
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) buildSharedWishlistMarkup(wishes []*entities.Wish, totalEntities int, offset int, sourceId string) *inlinekeyboard.Builder {
	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	keyboard.AppendAsLine(
		keyboard.NewButton(s.constants.BTN_REFRESH, s.builders.CallbackDataBuilder.Build(sourceId, s.constants.CMD_SHOW_SWL, strconv.Itoa(offset)).String()),
	).Append(
		s.BuildEntityButtons(wishes, offset, func(id string, offset int) string {
			return s.builders.CallbackDataBuilder.Build(id, s.constants.CMD_SHOW_SWI, strconv.Itoa(offset)).String()
		}),
	).Append(
		s.pagination.BuildControls(totalEntities, s.constants.CMD_SHOW_SWL, sourceId, offset),
	)

	return keyboard
}

func (s *Service) buildSharedWishInfoKeyboard(
	wish *entities.Wish,
	offset,
	sourceId string,
	viewerId string,
) *inlinekeyboard.Builder {
	keyboard := inlinekeyboard.New()

	if wish.Link != "" {
		keyboard.AppendAsLine(keyboard.NewURLButton(wish.GetMarketplace(), wish.Link))
	}

	if wish.ExecutorId != "" {
		if wish.ExecutorId == viewerId {
			keyboard.AppendAsLine(keyboard.NewButton(s.constants.BTN_CANCEL_BOOKING, s.builders.CallbackDataBuilder.Build(wish.ID, s.constants.CMD_TOGGLE_WISH_LOCK, offset).String()))
		}
		// has executor but it's not viewer - not show lock button
	} else {
		keyboard.AppendAsLine(keyboard.NewButton(s.constants.BTN_BOOK_WISH, s.builders.CallbackDataBuilder.Build(wish.ID, s.constants.CMD_TOGGLE_WISH_LOCK, offset).String()))
	}

	keyboard.AppendAsStack(
		keyboard.NewButton(s.constants.BTN_BACK_TO_WISHLIST, s.builders.CallbackDataBuilder.Build(sourceId, s.constants.CMD_SHOW_SWL, offset).String()),
	)

	return keyboard
}
