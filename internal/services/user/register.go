package user

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/entities"
)

func (s *Service) RegisterOrUpdateUser(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()

	userId := strconv.Itoa(message.From.Id)

	bf, err := entities.NewBaseFields(false, s.cfg.Timezone)
	if err != nil {
		return err
	}

	user := entities.User{
		BaseFields: bf,
		Name:       message.From.FirstName,
		TgUsername: message.From.Username,
		TgId:       userId,
		ChatId:     strconv.Itoa(message.Chat.Id),
		IsAdmin:    message.From.IsAdmin(),
	}

	err = s.repositories.User.Save(ctx, &user)
	if err != nil {
		return err
	}

	wishLists, err := s.repositories.WishList.Filter(ctx, &entities.WishList{UserId: userId})
	if err != nil {
		return err
	}

	bf, err = entities.NewBaseFields(true, s.cfg.Timezone)
	if err != nil {
		return err
	}

	if len(wishLists) == 0 {
		wishList := entities.WishList{
			BaseFields: bf,
			Name:       s.cfg.Constants.DEFAULT_WISHLIST_NAME,
			ChatId:     message.GetChatIdStr(),
			UserId:     userId,
		}
		err = s.repositories.WishList.Save(ctx, &wishList)
		if err != nil {
			s.clients.Telegram.Reply(ctx, s.cfg.Constants.WISHLIST_CREATION_ERROR, update)
			return err
		}
	}

	return nil
}
