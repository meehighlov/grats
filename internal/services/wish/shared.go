package wish

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/repositories/models"
	"github.com/meehighlov/grats/internal/repositories/wish"
	"github.com/meehighlov/grats/internal/repositories/wish_list"
	"github.com/meehighlov/grats/pkg/telegram"
	inlinekeyboard "github.com/meehighlov/grats/pkg/telegram/builders/inline_keyboard"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
)

func (s *Service) ShareWishList(ctx context.Context, scope *telegram.Scope) error {
	var (
		wishlist []*models.WishList
	)
	params := scope.CallbackData().FromString(scope.Update().CallbackQuery.Data)

	wishListId := params.ID

	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		wishlist, err = s.repositories.WishList.List(ctx, &wish_list.ListFilter{WishListID: wishListId})
		return err
	})
	if err != nil {
		return err
	}

	if len(wishlist) == 0 {
		return fmt.Errorf("wishlist not found")
	}

	botName := s.cfg.BotName
	shareLink := fmt.Sprintf(s.cfg.Constants.SHARE_WISHLIST_LINK_TEMPLATE, botName, wishlist[0].ID)

	keyboard := scope.Keyboard()
	keyboard.AppendAsLine(
		keyboard.NewShareLinkButton(s.cfg.Constants.BTN_SHARE, shareLink, s.cfg.Constants.MY_WISHLIST_SHARE_TITLE),
		keyboard.NewCopyButton(s.cfg.Constants.BTN_COPY_LINK, shareLink),
	)
	keyboard.AppendAsLine(
		keyboard.NewButton(s.cfg.Constants.BTN_BACK_TO_WISHLIST, scope.CallbackData().Build(wishlist[0].ID, s.cfg.Constants.CMD_LIST, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
	)

	shareMessage := s.cfg.Constants.SHARE_WISHLIST_MESSAGE

	if _, err := scope.Edit(
		ctx,
		shareMessage,
		tgc.WithReplyMurkup(keyboard.Murkup()),
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) ShowSharedWishlist(ctx context.Context, scope *telegram.Scope) error {
	var (
		wishes     []*models.Wish
		count      int64
		offset     int
		wishlistId string
		header     string
	)
	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		// case when called from /start or comes from link
		isFromStartOption := !scope.Update().IsCallback()

		if isFromStartOption {
			if err := s.userRegistration.RegisterOrUpdateUser(ctx, scope); err != nil {
				return err
			}
		}

		message := scope.Update().GetMessage()

		listPrefix := s.cfg.Constants.CMD_START + " " + s.cfg.Constants.SHARED_LIST_ID_PREFIX

		wishlistId = strings.TrimPrefix(message.Text, listPrefix)
		offset = s.cfg.Constants.LIST_START_OFFSET
		if scope.Update().IsCallback() {
			params := scope.CallbackData().FromString(scope.Update().CallbackQuery.Data)
			wishlistId = params.ID
			offset, _ = strconv.Atoi(params.Offset)
			if offset == 0 {
				offset = s.cfg.Constants.LIST_START_OFFSET
			}
		}

		wishes, err = s.repositories.Wish.List(ctx, &wish.ListFilter{WishListID: wishlistId, Limit: s.cfg.ListLimitLen, Offset: offset})
		if err != nil {
			scope.Reply(ctx, s.cfg.Constants.FAILED_TO_LOAD_WISHES)
			return err
		}

		if len(wishes) == 0 {
			// TODO: update message when user list is empty
			s.logger.Debug("Not found wishes to share")
			return nil
		}

		// user used his own link
		// just send greeting
		if s.cfg.IsProd() {
			if wishes[0].UserId == strconv.Itoa(message.From.Id) {
				scope.Reply(ctx, s.cfg.Constants.HELLO_AGAIN)
				return nil
			}
		}

		userInfo, _ := scope.GetChatMember(ctx, wishes[0].UserId)
		name := userInfo.Result.User.Username
		if name == "" {
			name = userInfo.Result.User.FirstName
		}
		header = fmt.Sprintf(s.cfg.Constants.WISHLIST_HEADER_TEMPLATE, "@"+name)

		count, err = s.repositories.Wish.Count(ctx, &wish.CountFilter{WishListID: wishlistId})
		return err
	})
	if err != nil {
		return err
	}

	keyboard := s.buildSharedWishlistMarkup(scope, wishes, int(count), offset, wishlistId)

	if scope.Update().IsCallback() {
		if _, err := scope.Edit(
			ctx,
			header,
			tgc.WithReplyMurkup(keyboard.Murkup()),
		); err != nil {
			return err
		}

		return nil
	}

	if _, err := scope.Reply(
		ctx,
		header,
		tgc.WithReplyMurkup(keyboard.Murkup()),
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) buildSharedWishlistMarkup(scope *telegram.Scope, wishes []*models.Wish, totalmodels int, offset int, sourceId string) *inlinekeyboard.Builder {
	keyboard := scope.Keyboard()

	keyboard.AppendAsLine(
		keyboard.NewButton(s.cfg.Constants.BTN_REFRESH, scope.CallbackData().Build(sourceId, s.cfg.Constants.CMD_SHOW_SWL, strconv.Itoa(offset)).String()),
	).Append(
		s.BuildEntityButtons(scope, wishes, offset, func(id string, offset int) string {
			return scope.CallbackData().Build(id, s.cfg.Constants.CMD_SHOW_SWI, strconv.Itoa(offset)).String()
		}),
	).Append(
		scope.Pagination().BuildControls(totalmodels, s.cfg.Constants.CMD_SHOW_SWL, sourceId, offset),
	)

	return keyboard
}

func (s *Service) buildSharedWishInfoKeyboard(
	scope *telegram.Scope,
	wish *models.Wish,
	offset,
	sourceId string,
	viewerId string,
) *inlinekeyboard.Builder {
	keyboard := scope.Keyboard()

	if wish.Link != "" {
		keyboard.AppendAsLine(scope.Keyboard().NewURLButton(wish.GetMarketplace(s.GetSiteName), wish.Link))
	}

	if wish.ExecutorId != "" {
		if wish.ExecutorId == viewerId {
			keyboard.AppendAsLine(keyboard.NewButton(s.cfg.Constants.BTN_CANCEL_BOOKING, scope.CallbackData().Build(wish.ID, s.cfg.Constants.CMD_TOGGLE_WISH_LOCK, offset).String()))
		}
		// has executor but it's not viewer - not show lock button
	} else {
		keyboard.AppendAsLine(keyboard.NewButton(s.cfg.Constants.BTN_BOOK_WISH, scope.CallbackData().Build(wish.ID, s.cfg.Constants.CMD_TOGGLE_WISH_LOCK, offset).String()))
	}

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_BACK_TO_WISHLIST, scope.CallbackData().Build(sourceId, s.cfg.Constants.CMD_SHOW_SWL, offset).String()),
	)

	return keyboard
}
