package wish

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/repositories/models"
	"github.com/meehighlov/grats/internal/repositories/wish"
	"github.com/meehighlov/grats/pkg/telegram"
	inlinekeyboard "github.com/meehighlov/grats/pkg/telegram/builders/inline_keyboard"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
)

func (s *Service) List(ctx context.Context, scope *telegram.Scope) error {
	var (
		listId  string
		userId  string
		offset  string = s.cfg.Constants.LIST_DEFAULT_OFFSET
		wishes  []*models.Wish
		count   int64
		offset_ int
	)
	err := s.db.Tx(ctx, func(ctx context.Context) (err error) {
		userId = strconv.Itoa(scope.Update().GetMessage().From.Id)
		if scope.Update().IsCallback() {
			callbackData := scope.CallbackData().FromString(scope.Update().CallbackQuery.Data)
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

	if scope.Update().IsCallback() {
		if _, err := scope.Edit(
			ctx,
			s.buildListHeaderMessage(wishes),
			tgc.WithReplyMurkup(s.buildListMarkup(scope, int(count), wishes, offset_, listId).Murkup()),
		); err != nil {
			return err
		}
	} else {
		if _, err := scope.Reply(
			ctx,
			s.buildListHeaderMessage(wishes),
			tgc.WithReplyMurkup(s.buildListMarkup(scope, int(count), wishes, offset_, listId).Murkup()),
		); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) buildListMarkup(scope *telegram.Scope, totalmodels int, models []*models.Wish, offset int, listId string) *inlinekeyboard.Builder {
	callbackData := scope.CallbackData()
	callbackDataBuilder := func(id string, offset int) string {
		return callbackData.Build(id, s.cfg.Constants.CMD_WISH_INFO, strconv.Itoa(offset)).String()
	}
	entityListAsButtons := s.BuildEntityButtons(scope, models, offset, callbackDataBuilder)
	keyboard := scope.Keyboard()

	if len(models) > 0 {
		keyboard.AppendAsLine(
			keyboard.NewButton(s.cfg.Constants.BTN_ADD_WISH, callbackData.Build(listId, s.cfg.Constants.CMD_ADD_TO_WISH, strconv.Itoa(offset)).String()),
			keyboard.NewButton(s.cfg.Constants.BTN_SHARE_LIST, callbackData.Build(listId, s.cfg.Constants.CMD_SHARE_WISH_LIST, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
		)
	} else {
		keyboard.AppendAsLine(
			keyboard.NewButton(s.cfg.Constants.BTN_ADD_WISH, callbackData.Build(listId, s.cfg.Constants.CMD_ADD_TO_WISH, strconv.Itoa(offset)).String()),
		)
	}

	controls := scope.Pagination().BuildControls(totalmodels, s.cfg.Constants.CMD_LIST, listId, offset)

	return keyboard.Append(entityListAsButtons).Append(controls)
}

func (s *Service) buildListHeaderMessage(wishes []*models.Wish) string {
	if len(wishes) == 0 {
		return s.cfg.Constants.WISHLIST_EMPTY
	}
	return s.cfg.Constants.MY_WISHLIST
}
