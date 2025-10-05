package wish

import (
	"context"
	"strconv"

	inlinekeyboard "github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/entities"
	"github.com/meehighlov/grats/internal/repositories/wish"
)

func (s *Service) List(ctx context.Context, update *telegram.Update) error {
	var (
		listId string
		userId string
		offset string = s.constants.LIST_DEFAULT_OFFSET
	)
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

	offset_, _ := strconv.Atoi(offset)
	if offset_ == 0 {
		offset_ = s.constants.LIST_START_OFFSET
	}

	entities, err := s.repositories.Wish.List(ctx, &wish.ListFilter{WishListID: listId, Limit: s.cfg.ListLimitLen, Offset: offset_})
	if err != nil {
		return err
	}

	count, err := s.repositories.Wish.Count(ctx, &wish.CountFilter{WishListID: listId})
	if err != nil {
		return err
	}

	if update.IsCallback() {
		if _, err := s.clients.Telegram.Edit(
			ctx,
			s.buildListHeaderMessage(entities),
			update,
			telegram.WithReplyMurkup(s.buildListMarkup(int(count), entities, offset_, listId).Murkup()),
		); err != nil {
			return err
		}
	} else {
		if _, err := s.clients.Telegram.Reply(
			ctx,
			s.buildListHeaderMessage(entities),
			update,
			telegram.WithReplyMurkup(s.buildListMarkup(int(count), entities, offset_, listId).Murkup()),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) buildListMarkup(totalEntities int, entities []*entities.Wish, offset int, listId string) *inlinekeyboard.Builder {
	callbackDataBuilder := func(id string, offset int) string {
		return s.builders.CallbackDataBuilder.Build(id, s.constants.CMD_WISH_INFO, strconv.Itoa(offset)).String()
	}
	entityListAsButtons := s.BuildEntityButtons(entities, offset, callbackDataBuilder)
	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	if len(entities) > 0 {
		keyboard.AppendAsLine(
			s.builders.KeyboardBuilder.NewButton(s.constants.BTN_ADD_WISH, s.builders.CallbackDataBuilder.Build(listId, s.constants.CMD_ADD_TO_WISH, strconv.Itoa(offset)).String()),
			s.builders.KeyboardBuilder.NewButton(s.constants.BTN_SHARE_LIST, s.builders.CallbackDataBuilder.Build(listId, s.constants.CMD_SHARE_WISH_LIST, s.constants.LIST_DEFAULT_OFFSET).String()),
		)
	} else {
		keyboard.AppendAsLine(
			s.builders.KeyboardBuilder.NewButton(s.constants.BTN_ADD_WISH, s.builders.CallbackDataBuilder.Build(listId, s.constants.CMD_ADD_TO_WISH, strconv.Itoa(offset)).String()),
		)
	}

	controls := s.pagination.BuildControls(totalEntities, s.constants.CMD_LIST, listId, offset)

	return keyboard.Append(entityListAsButtons).Append(controls)
}

func (s *Service) buildListHeaderMessage(wishes []*entities.Wish) string {
	if len(wishes) == 0 {
		return s.constants.WISHLIST_EMPTY
	}
	return s.constants.MY_WISHLIST
}
