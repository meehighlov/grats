package wish

import (
	"context"
	"errors"
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"

	inlinekeyboard "github.com/meehighlov/grats/internal/builders/inline_keyboard"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/entities"
)

func (s *Service) AddWish(ctx context.Context, update *telegram.Update) error {
	var (
		wishListId string
		userId     string
	)

	userId = strconv.Itoa(update.GetMessage().From.Id)
	wishListId = s.builders.CallbackDataBuilder.FromString(update.CallbackQuery.Data).ID

	wishes, err := s.repositories.Wish.Filter(ctx, &entities.Wish{UserId: userId})
	if err != nil {
		s.clients.Telegram.Reply(ctx, s.cfg.Constants.ERROR_MESSAGE, update)
		return err
	}

	if len(wishes) >= s.cfg.Constants.WISH_LIMIT_FOR_USER {
		s.clients.Telegram.Reply(ctx, fmt.Sprintf(
			s.cfg.Constants.WISH_LIMIT_REACHED_TEMPLATE,
			s.cfg.Constants.WISH_LIMIT_FOR_USER,
		), update)
		return nil
	}

	msg := s.cfg.Constants.ENTER_WISH_NAME

	if _, err := s.clients.Telegram.Reply(ctx, msg, update); err != nil {
		return err
	}
	s.clients.Cache.AppendText(ctx, update.GetChatIdStr(), wishListId)
	s.clients.Cache.AppendText(ctx, update.GetChatIdStr(), userId)

	return nil
}

func (s *Service) SaveWish(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()
	userId := strconv.Itoa(message.From.Id)
	texts, err := s.clients.Cache.GetTexts(ctx, update.GetChatIdStr())
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
		s.clients.Telegram.Reply(ctx, fmt.Sprintf(s.cfg.Constants.WISH_NAME_TOO_LONG_TEMPLATE, s.cfg.Constants.WISH_NAME_MAX_LEN), update)
		return errors.New("wish name is too long")
	}

	validatedName, err := s.validateWishName(message.Text)
	if err != nil {
		s.clients.Telegram.Reply(ctx, s.cfg.Constants.WISH_NAME_INVALID_CHARS, update)
		return errors.New("wish name contains invalid characters")
	}

	bf, err := entities.NewBaseFields(false, s.cfg.Timezone)
	if err != nil {
		return err
	}

	wish := entities.Wish{
		BaseFields: bf,
		Name:       validatedName,
		ChatId:     message.GetChatIdStr(),
		UserId:     userId,
		WishListId: wishListId,
	}

	err = s.repositories.Wish.Save(ctx, &wish)
	if err != nil {
		return err
	}

	msg := s.cfg.Constants.WISH_ADDED

	if _, err := s.clients.Telegram.Reply(ctx, msg, update, telegram.WithReplyMurkup(s.buildWishNavigationMarkup(wish.ID, wishListId).Murkup())); err != nil {
		return err
	}

	return nil
}

func (s *Service) buildWishNavigationMarkup(wishId, wishListId string) *inlinekeyboard.Builder {
	keyboard := s.builders.KeyboardBuilder.NewKeyboard()

	keyboard.AppendAsStack(
		keyboard.NewButton(s.cfg.Constants.BTN_OPEN_WISH, s.builders.CallbackDataBuilder.Build(wishId, s.cfg.Constants.CMD_WISH_INFO, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_NEW_WISH, s.builders.CallbackDataBuilder.Build(wishListId, s.cfg.Constants.CMD_ADD_TO_WISH, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
		keyboard.NewButton(s.cfg.Constants.BTN_WISH_LIST, s.builders.CallbackDataBuilder.Build(wishListId, s.cfg.Constants.CMD_LIST, s.cfg.Constants.LIST_DEFAULT_OFFSET).String()),
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
