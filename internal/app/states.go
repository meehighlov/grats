package app

import (
	"context"

	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/repositories"
	"github.com/meehighlov/grats/internal/services"
	"github.com/meehighlov/grats/pkg/telegram"
)

func RegisterHandlers(
	b *telegram.Bot,
	s *services.Services,
	cfg *config.Config,
	repositories *repositories.Repositories,
) {
	c := cfg.Constants

	resetUserCache := func(ctx context.Context, scope *telegram.Scope) error {
		return repositories.Cache.Reset(ctx, scope.Update().GetChatIdStr())
	}

	b.AddMiddleware(
		func(ctx context.Context, scope *telegram.Scope) error {
			return scope.AnswerCallbackQuery(ctx)
		},
	)

	b.AddHandler(
		s.User.Start,
		telegram.Command(c.CMD_START),
	)

	// ------------------------ support --------------------------------

	b.AddHandler(
		s.Support.Support,
		telegram.Command(c.CMD_SUPPORT),
	)

	supportWrite := b.AddHandler(
		s.Support.SupportWrite,
		telegram.CallbackDataContains(c.CMD_SUPPORT_WRITE),
		telegram.BeforeAction(resetUserCache),
	)

	b.AddHandler(
		s.Support.SendSupportMessage,
		telegram.MessageHasText(),
		telegram.AcceptFrom(supportWrite),
	)

	b.AddHandler(
		s.Support.CancelSupportCall,
		telegram.CallbackDataContains(c.CMD_SUPPORT_CANCEL),
	)

	b.AddHandler(
		s.Support.ProcessSupportReply,
		SupportReplyCondition(cfg.SupportChatId),
	)

	// ------------------------ wishlist --------------------------------

	addWish := b.AddHandler(
		s.Wish.AddWish,
		telegram.CallbackDataContains(c.CMD_ADD_TO_WISH),
		telegram.BeforeAction(resetUserCache),
	)

	b.AddHandler(
		s.Wish.SaveWish,
		telegram.MessageHasText(),
		telegram.AcceptFrom(addWish),
	)

	b.AddHandler(
		s.Wish.List,
		telegram.Command(c.CMD_WISHLIST),
	)

	b.AddHandler(
		s.Wish.List,
		telegram.CallbackDataContains(c.CMD_LIST),
	)

	b.AddHandler(
		s.Wish.WishInfo,
		telegram.CallbackDataContains(c.CMD_WISH_INFO),
	)

	b.AddHandler(
		s.Wish.DeleteWish,
		telegram.CallbackDataContains(c.CMD_DELETE_WISH),
	)

	b.AddHandler(
		s.Wish.ConfirmDeleteWish,
		telegram.CallbackDataContains(c.CMD_CONFIRM_DELETE_WISH),
	)

	editWishName := b.AddHandler(
		s.Wish.EditWishName,
		telegram.CallbackDataContains(c.CMD_EDIT_WISH_NAME),
		telegram.BeforeAction(resetUserCache),
	)

	b.AddHandler(
		s.Wish.SaveEditWishName,
		telegram.MessageHasText(),
		telegram.AcceptFrom(editWishName),
	)

	editWishLink := b.AddHandler(
		s.Wish.EditLink,
		telegram.CallbackDataContains(c.CMD_EDIT_LINK),
		telegram.BeforeAction(resetUserCache),
	)

	b.AddHandler(
		s.Wish.SaveEditLink,
		telegram.MessageHasText(),
		telegram.AcceptFrom(editWishLink),
	)

	editWishPrice := b.AddHandler(
		s.Wish.EditPrice,
		telegram.CallbackDataContains(c.CMD_EDIT_PRICE),
		telegram.BeforeAction(resetUserCache),
	)

	b.AddHandler(
		s.Wish.SaveEditPrice,
		telegram.MessageHasText(),
		telegram.AcceptFrom(editWishPrice),
	)

	b.AddHandler(
		s.Wish.DeleteLink,
		telegram.CallbackDataContains(c.CMD_DELETE_LINK),
		telegram.AcceptFrom(editWishLink),
	)

	b.AddHandler(
		s.Wish.ShareWishList,
		telegram.CallbackDataContains(c.CMD_SHARE_WISH_LIST),
	)

	b.AddHandler(
		s.Wish.ToggleWishLock,
		telegram.CallbackDataContains(c.CMD_TOGGLE_WISH_LOCK),
	)

	b.AddHandler(
		s.Wish.ShowSharedWishlist,
		ShowSharedListCondition(),
	)

	b.AddHandler(
		s.Wish.WishInfo,
		telegram.CallbackDataContains(c.CMD_SHOW_SWI),
	)

	b.AddHandler(
		s.Wish.ShowSharedWishlist,
		telegram.CallbackDataContains(c.CMD_SHOW_SWL),
	)

	b.Reset(
		telegram.Command(c.CMD_CANCEL),
		telegram.BeforeAction(resetUserCache),
		telegram.AcceptFrom(supportWrite),
		telegram.AcceptFrom(editWishName),
		telegram.AcceptFrom(editWishLink),
		telegram.AcceptFrom(editWishPrice),
		telegram.AcceptFrom(addWish),
	)
}
