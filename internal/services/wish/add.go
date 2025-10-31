package wish

import (
	"context"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/repositories/models"
	"github.com/meehighlov/grats/internal/repositories/wish"
	"github.com/meehighlov/grats/pkg/telegram"
	inlinekeyboard "github.com/meehighlov/grats/pkg/telegram/builders/inline_keyboard"
	tgc "github.com/meehighlov/grats/pkg/telegram/client"
)

func (s *Service) AddWish(ctx context.Context, scope *telegram.Scope) error {
	var (
		wishListId string
		userId     string
		wishes     []*models.Wish
	)

	userId = strconv.Itoa(scope.Update().GetMessage().From.Id)
	wishListId = scope.CallbackData().FromString(scope.Update().CallbackQuery.Data).ID

	err := s.db.Tx(ctx, func(ctx context.Context) error {
		var err error
		wishes, err = s.repositories.Wish.List(ctx, &wish.ListFilter{UserId: userId})
		if err != nil {
			scope.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE)
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(wishes) >= s.cfg.Constants.WISH_LIMIT_FOR_USER {
		scope.Reply(ctx, fmt.Sprintf(
			s.cfg.Constants.WISH_LIMIT_REACHED_TEMPLATE,
			s.cfg.Constants.WISH_LIMIT_FOR_USER,
		))
		return nil
	}

	msg := s.cfg.Constants.ENTER_WISH_NAME

	if _, err := scope.Reply(ctx, msg); err != nil {
		return err
	}
	s.repositories.Cache.AppendText(ctx, scope.Update().GetChatIdStr(), wishListId)
	s.repositories.Cache.AppendText(ctx, scope.Update().GetChatIdStr(), userId)

	return nil
}

func (s *Service) SaveWish(ctx context.Context, scope *telegram.Scope) error {
	message := scope.Update().GetMessage()
	userId := strconv.Itoa(message.From.Id)
	texts, err := s.repositories.Cache.GetTexts(ctx, scope.Update().GetChatIdStr())
	if err != nil {
		return err
	}

	wishListId := texts[0]

	if wishListId == "" {
		s.logger.Error(
			"SaveWish",
			"error", "wishListId is empty",
			"userId", userId,
		)
		return nil
	}

	if len(message.Text) > s.cfg.Constants.WISH_NAME_MAX_LEN {
		scope.Reply(ctx, fmt.Sprintf(s.cfg.Constants.WISH_NAME_TOO_LONG_TEMPLATE, s.cfg.Constants.WISH_NAME_MAX_LEN))
		return errors.New("wish name is too long")
	}

	validatedName, err := s.validateWishName(message.Text)
	if err != nil {
		scope.Reply(ctx, s.cfg.Constants.WISH_NAME_INVALID_CHARS)
		return errors.New("wish name contains invalid characters")
	}

	bf, err := models.NewBaseFields(false, s.cfg.Timezone)
	if err != nil {
		return err
	}

	wish := models.Wish{
		BaseFields: bf,
		Name:       validatedName,
		ChatId:     message.GetChatIdStr(),
		UserId:     userId,
		WishListId: wishListId,
	}

	err = s.db.Tx(ctx, func(ctx context.Context) error {
		return s.repositories.Wish.Save(ctx, &wish)
	})
	if err != nil {
		return err
	}

	msg := s.cfg.Constants.WISH_ADDED

	if _, err := scope.Reply(ctx, msg, tgc.WithReplyMurkup(s.buildWishNavigationMarkup(scope, wish.ID, wishListId).Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) buildWishNavigationMarkup(scope *telegram.Scope, wishId, wishListId string) *inlinekeyboard.Builder {
	keyboard := scope.Keyboard()

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_OPEN_WISH, scope.CallbackData().Build(wishId, s.cfg.Constants.CMD_WISH_INFO, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_NEW_WISH, scope.CallbackData().Build(wishListId, s.cfg.Constants.CMD_ADD_TO_WISH, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_WISH_LIST, scope.CallbackData().Build(wishListId, s.cfg.Constants.CMD_LIST, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
	)

	return keyboard
}

func (s *Service) validateWishName(name string) (string, error) {
	trimmedName := strings.TrimSpace(name)
	if len(trimmedName) > s.cfg.Constants.WISH_NAME_MAX_LEN {
		return "", fmt.Errorf("название желания не должно превышать %d символов", s.cfg.Constants.WISH_NAME_MAX_LEN)
	}

	sanitizedName := html.EscapeString(trimmedName)

	reg := regexp.MustCompile(`[^\p{L}\p{N}\p{P}\p{Z}]`)
	sanitizedName = reg.ReplaceAllString(sanitizedName, "")

	return sanitizedName, nil
}
