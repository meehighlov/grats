package handlers

import (
	"context"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

func RegisterOrUpdateUser(ctx context.Context, event *common.Event) error {
	err := db.GetDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		message := event.GetMessage()

		isAdmin := 0
		if message.From.IsAdmin() {
			isAdmin = 1
		}

		userId := strconv.Itoa(message.From.Id)

		user := db.User{
			BaseFields: db.NewBaseFields(false),
			Name:       message.From.FirstName,
			TgUsername: message.From.Username,
			TgId:       userId,
			ChatId:     strconv.Itoa(message.Chat.Id),
			Birthday:   "",
			IsAdmin:    isAdmin,
		}

		err := user.Save(ctx, tx)
		if err != nil {
			event.Logger.Error("start error creating user", "chatId", message.GetChatIdStr(), "error", err.Error())
			return err
		}

		chat := db.Chat{
			BaseFields:     db.NewBaseFields(false),
			ChatType:       "private",
			ChatId:         event.GetMessage().GetChatIdStr(),
			BotInvitedById: userId,
			GreetingTemplate: "🔔Сегодня день рождения празднует %s🥳",
		}

		err = chat.Save(ctx, tx)
		if err != nil {
			event.Logger.Error("start error creating chat", "chatId", message.GetChatIdStr(), "error", err.Error())
			return err
		}

		wishLists, err := (&db.WishList{UserId: userId}).Filter(ctx, tx)
		if err != nil {
			event.Logger.Error("start error getting wishlists", "chatId", message.GetChatIdStr(), "error", err.Error())
			return err
		}

		if len(wishLists) == 0 {
			wishList := db.WishList{
				BaseFields: db.NewBaseFields(true),
				Name:       "Мой wishlist",
				ChatId:     message.GetChatIdStr(),
				UserId:     userId,
			}

			err = wishList.Save(ctx, tx)
			if err != nil {
				event.Logger.Error("error creating wish_list: " + err.Error())
				event.Reply(ctx, "Возникла непредвиденная ошибка при создании первого списка желаний, над этим уже работают😔")
				return err
			}
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
