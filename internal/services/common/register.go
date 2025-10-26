package common

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/repositories/models"
	"github.com/meehighlov/grats/internal/repositories/wish_list"
)

func (u *UserRegistration) RegisterOrUpdateUser(ctx context.Context, update *telegram.Update) error {
	message := update.GetMessage()

	userId := strconv.Itoa(message.From.Id)

	bf, err := models.NewBaseFields(false, u.cfg.Timezone)
	if err != nil {
		return err
	}

	user := models.User{
		BaseFields: bf,
		Name:       message.From.FirstName,
		TgUsername: message.From.Username,
		TgId:       userId,
		ChatId:     strconv.Itoa(message.Chat.Id),
		IsAdmin:    message.From.IsAdmin(),
	}

	err = u.repositories.User.Save(ctx, &user)
	if err != nil {
		return err
	}

	wishLists, err := u.repositories.WishList.List(ctx, &wish_list.ListFilter{UserId: userId})
	if err != nil {
		return err
	}

	bf, err = models.NewBaseFields(true, u.cfg.Timezone)
	if err != nil {
		return err
	}

	if len(wishLists) == 0 {
		wishList := models.WishList{
			BaseFields: bf,
			Name:       u.cfg.Constants.DEFAULT_WISHLIST_NAME,
			ChatId:     message.GetChatIdStr(),
			UserId:     userId,
		}
		err = u.repositories.WishList.Save(ctx, &wishList)
		if err != nil {
			u.clients.Telegram.Reply(ctx, u.cfg.Constants.WISHLIST_CREATION_ERROR, update)
			return err
		}
	}

	return nil
}
