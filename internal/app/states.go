package app

import (
	"context"

	appconditions "github.com/meehighlov/grats/internal/app/conditions"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/infra/postgres"
	"github.com/meehighlov/grats/internal/repositories"
	"github.com/meehighlov/grats/internal/services"
	tfsm "github.com/meehighlov/grats/pkg/telegram/fsm"
	"github.com/meehighlov/grats/pkg/telegram/fsm/when"
	"github.com/meehighlov/grats/pkg/telegram/fsm/with"
	tgm "github.com/meehighlov/grats/pkg/telegram/models"
)

func RegisterStates(
	f *tfsm.FSM,
	s *services.Services,
	cfg *config.Config,
	clients *clients.Clients,
	repositories *repositories.Repositories,
	db *postgres.DB,
) {
	c := cfg.Constants

	resetUserCache := func(ctx context.Context, update *tgm.Update) error {
		return repositories.Cache.Reset(ctx, update.GetChatIdStr())
	}

	f.AddMiddleware(clients.Telegram.AnswerCallbackQuery)

	f.Activate(
		s.User.Start,
		when.Command(c.CMD_START),
	)

	// ------------------------ support --------------------------------

	f.Activate(
		s.Support.Support,
		when.Command(c.CMD_SUPPORT),
	)

	supportWrite := f.Activate(
		s.Support.SupportWrite,
		when.CallbackDataContains(c.CMD_SUPPORT_WRITE),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		s.Support.SendSupportMessage,
		when.MessageHasText(),
		with.AcceptFrom(supportWrite),
	)

	f.Activate(
		s.Support.CancelSupportCall,
		when.CallbackDataContains(c.CMD_SUPPORT_CANCEL),
	)

	f.Activate(
		s.Support.ProcessSupportReply,
		appconditions.SupportReplyCondition(cfg.SupportChatId),
	)

	// ------------------------ wishlist --------------------------------

	addWish := f.Activate(
		s.Wish.AddWish,
		when.CallbackDataContains(c.CMD_ADD_TO_WISH),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		s.Wish.SaveWish,
		when.MessageHasText(),
		with.AcceptFrom(addWish),
	)

	f.Activate(
		s.Wish.List,
		when.Command(c.CMD_WISHLIST),
	)

	f.Activate(
		s.Wish.List,
		when.CallbackDataContains(c.CMD_LIST),
	)

	f.Activate(
		s.Wish.WishInfo,
		when.CallbackDataContains(c.CMD_WISH_INFO),
	)

	f.Activate(
		s.Wish.DeleteWish,
		when.CallbackDataContains(c.CMD_DELETE_WISH),
	)

	f.Activate(
		s.Wish.ConfirmDeleteWish,
		when.CallbackDataContains(c.CMD_CONFIRM_DELETE_WISH),
	)

	editWishName := f.Activate(
		s.Wish.EditWishName,
		when.CallbackDataContains(c.CMD_EDIT_WISH_NAME),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		s.Wish.SaveEditWishName,
		when.MessageHasText(),
		with.AcceptFrom(editWishName),
	)

	editWishLink := f.Activate(
		s.Wish.EditLink,
		when.CallbackDataContains(c.CMD_EDIT_LINK),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		s.Wish.SaveEditLink,
		when.MessageHasText(),
		with.AcceptFrom(editWishLink),
	)

	editWishPrice := f.Activate(
		s.Wish.EditPrice,
		when.CallbackDataContains(c.CMD_EDIT_PRICE),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		s.Wish.SaveEditPrice,
		when.MessageHasText(),
		with.AcceptFrom(editWishPrice),
	)

	f.Activate(
		s.Wish.DeleteLink,
		when.CallbackDataContains(c.CMD_DELETE_LINK),
		with.AcceptFrom(editWishLink),
	)

	f.Activate(
		s.Wish.ShareWishList,
		when.CallbackDataContains(c.CMD_SHARE_WISH_LIST),
	)

	f.Activate(
		s.Wish.ToggleWishLock,
		when.CallbackDataContains(c.CMD_TOGGLE_WISH_LOCK),
	)

	f.Activate(
		s.Wish.ShowSharedWishlist,
		appconditions.ShowSharedListCondition(),
	)

	f.Activate(
		s.Wish.WishInfo,
		when.CallbackDataContains(c.CMD_SHOW_SWI),
	)

	f.Activate(
		s.Wish.ShowSharedWishlist,
		when.CallbackDataContains(c.CMD_SHOW_SWL),
	)

	f.Reset(
		when.Command(c.CMD_CANCEL),
		with.BeforeAction(resetUserCache),
		with.AcceptFrom(supportWrite),
		with.AcceptFrom(editWishName),
		with.AcceptFrom(editWishLink),
		with.AcceptFrom(editWishPrice),
		with.AcceptFrom(addWish),
	)
}
