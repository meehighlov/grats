package wish

import (
	"context"
	"strconv"

	inlinekeyboard "github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/models"
	"github.com/meehighlov/grats/internal/repositories/wish"
)

func (s *Service) List(ctx context.Context, update *telegram.Update) error {
	var (
		listId  string
		userId  string
		offset  string = s.cfg.Constants.LIST_DEFAULT_OFFSET
		wishes  []*models.Wish
		count   int64
		offset_ int
	)
	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		userId = strconv.Itoa(update.GetMessage().From.Id)
		if update.IsCallback() {
			callbackData := s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data)
			listId = callbackData.ID
			offset = callbackData.Offset
		} else {
			l, err := s.PickFirstWishList(ctx, userId)
			if err != nil {
				return err
			}
			listId = l.ID
		}

		offset_, err = strconv.Atoi(offset)
		if err != nil || offset_ == 0 {
			offset_ = s.cfg.Constants.LIST_START_OFFSET
		}

		wishes, err = s.repositories.Wish.List(ctx, &wish.ListFilter{WishListID: listId, Limit: s.cfg.ListLimitLen, Offset: offset_})
		if err != nil {
			return err
		}

		count, err = s.repositories.Wish.Count(ctx, &wish.CountFilter{WishListID: listId})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	if update.IsCallback() {
		if _, err := s.clients.Telegram.Edit(
			ctx,
			s.buildListHeaderMessage(wishes),
			update,
			telegram.WithReplyMurkup(s.buildListMarkup(int(count), wishes, offset_, listId).Murkup()),
		); err != nil {
			return err
		}
	} else {
		if _, err := s.clients.Telegram.Reply(
			ctx,
			s.buildListHeaderMessage(wishes),
			update,
			telegram.WithReplyMurkup(s.buildListMarkup(int(count), wishes, offset_, listId).Murkup()),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) buildListMarkup(totalmodels int, models []*models.Wish, offset int, listId string) *inlinekeyboard.Builder {
	callbackDataBuilder := func(id string, offset int) string {
		return s.builders.CallbackDataBuilder.Build(id, s.cfg.Constants.CMD_WISH_INFO, strconv.Itoa(offset)).String()
	}
	entityListAsButtons := s.BuildEntityButtons(models, offset, callbackDataBuilder)
	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	if len(models) > 0 {
		keyboard.AppendAsLine(
			s.builders.KeyboardBuilder.NewButton(s.cfg.Constants.BTN_ADD_WISH, s.builders.CallbackDataBuilder.Build(listId, s.cfg.Constants.CMD_ADD_TO_WISH, strconv.Itoa(offset)).String()),
			s.builders.KeyboardBuilder.NewButton(s.cfg.Constants.BTN_SHARE_LIST, s.builders.CallbackDataBuilder.Build(listId, s.cfg.Constants.CMD_SHARE_WISH_LIST, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
		)
	} else {
		keyboard.AppendAsLine(
			s.builders.KeyboardBuilder.NewButton(s.cfg.Constants.BTN_ADD_WISH, s.builders.CallbackDataBuilder.Build(listId, s.cfg.Constants.CMD_ADD_TO_WISH, strconv.Itoa(offset)).String()),
		)
	}

	controls := s.builders.PaginationBuilder.BuildControls(totalmodels, s.cfg.Constants.CMD_LIST, listId, offset)

	return keyboard.Append(entityListAsButtons).Append(controls)
}

func (s *Service) buildListHeaderMessage(wishes []*models.Wish) string {
	if len(wishes) == 0 {
		return s.cfg.Constants.WISHLIST_EMPTY
	}
	return s.cfg.Constants.MY_WISHLIST
}
