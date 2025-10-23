package app

import (
	"context"

	appconditions "github.com/meehighlov/grats/internal/app/conditions"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/db"
	"github.com/meehighlov/grats/internal/fsm"
	"github.com/meehighlov/grats/internal/fsm/when"
	"github.com/meehighlov/grats/internal/fsm/with"
	"github.com/meehighlov/grats/internal/services"
)

func RegisterStates(
	f *fsm.FSM,
	s *services.Services,
	cfg *config.Config,
	clients *clients.Clients,
	t *db.Tx,
) {
	c := cfg.Constants

	resetUserCache := func(ctx context.Context, update *telegram.Update) error {
		return clients.Cache.Reset(ctx, update.GetChatIdStr())
	}

	f.AddMiddleware(clients.Telegram.AnswerCallbackQuery)

	f.Activate(
		t.Wrap(s.User.Start),
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
		t.Wrap(s.Wish.AddWish),
		when.CallbackDataContains(c.CMD_ADD_TO_WISH),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		t.Wrap(s.Wish.SaveWish),
		when.MessageHasText(),
		with.AcceptFrom(addWish),
	)

	f.Activate(
		t.Wrap(s.Wish.List),
		when.Command(c.CMD_WISHLIST),
	)

	f.Activate(
		t.Wrap(s.Wish.List),
		when.CallbackDataContains(c.CMD_LIST),
	)

	f.Activate(
		t.Wrap(s.Wish.WishInfo),
		when.CallbackDataContains(c.CMD_WISH_INFO),
	)

	f.Activate(
		t.Wrap(s.Wish.DeleteWish),
		when.CallbackDataContains(c.CMD_DELETE_WISH),
	)

	f.Activate(
		t.Wrap(s.Wish.ConfirmDeleteWish),
		when.CallbackDataContains(c.CMD_CONFIRM_DELETE_WISH),
	)

	editWishName := f.Activate(
		t.Wrap(s.Wish.EditWishName),
		when.CallbackDataContains(c.CMD_EDIT_WISH_NAME),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		t.Wrap(s.Wish.SaveEditWishName),
		when.MessageHasText(),
		with.AcceptFrom(editWishName),
	)

	editWishLink := f.Activate(
		t.Wrap(s.Wish.EditLink),
		when.CallbackDataContains(c.CMD_EDIT_LINK),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		t.Wrap(s.Wish.SaveEditLink),
		when.MessageHasText(),
		with.AcceptFrom(editWishLink),
	)

	editWishPrice := f.Activate(
		t.Wrap(s.Wish.EditPrice),
		when.CallbackDataContains(c.CMD_EDIT_PRICE),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		t.Wrap(s.Wish.SaveEditPrice),
		when.MessageHasText(),
		with.AcceptFrom(editWishPrice),
	)

	f.Activate(
		t.Wrap(s.Wish.DeleteLink),
		when.CallbackDataContains(c.CMD_DELETE_LINK),
		with.AcceptFrom(editWishLink),
	)

	f.Activate(
		t.Wrap(s.Wish.ShareWishList),
		when.CallbackDataContains(c.CMD_SHARE_WISH_LIST),
	)

	f.Activate(
		t.Wrap(s.Wish.ToggleWishLock),
		when.CallbackDataContains(c.CMD_TOGGLE_WISH_LOCK),
	)

	f.Activate(
		t.Wrap(s.Wish.ShowSharedWishlist),
		appconditions.ShowSharedListCondition(),
	)

	f.Activate(
		t.Wrap(s.Wish.WishInfo),
		when.CallbackDataContains(c.CMD_SHOW_SWI),
	)

	f.Activate(
		t.Wrap(s.Wish.ShowSharedWishlist),
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
