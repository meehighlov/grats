package app

import (
	"context"

	appconditions "github.com/meehighlov/grats/internal/app/conditions"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/constants"
	"github.com/meehighlov/grats/internal/fsm"
	"github.com/meehighlov/grats/internal/fsm/when"
	"github.com/meehighlov/grats/internal/fsm/with"
	"github.com/meehighlov/grats/internal/orchestrators"
)

func RegisterStates(
	f *fsm.FSM,
	o *orchestrators.Orchestrators,
	c *constants.Constants,
	cfg *config.Config,
	clients *clients.Clients,
) {
	resetUserCache := func(ctx context.Context, update *telegram.Update) error {
		return clients.Cache.Reset(ctx, update.GetChatIdStr())
	}

	f.AddMiddleware(clients.Telegram.AnswerCallbackQuery)

	f.Activate(
		o.User.Start,
		when.Command(c.CMD_START),
	)

	// ------------------------ support --------------------------------

	f.Activate(
		o.Support.Support,
		when.Command(c.CMD_SUPPORT),
	)

	supportWrite := f.Activate(
		o.Support.SupportWrite,
		when.CallbackDataContains(c.CMD_SUPPORT_WRITE),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Support.SendSupportMessage,
		when.MessageHasText(),
		with.AcceptFrom(supportWrite),
	)

	f.Activate(
		o.Support.CancelSupportCall,
		when.CallbackDataContains(c.CMD_SUPPORT_CANCEL),
	)

	f.Activate(
		o.Support.ProcessSupportReply,
		appconditions.SupportReplyCondition(cfg.SupportChatId),
	)

	// ------------------------ wishlist --------------------------------

	addWish := f.Activate(
		o.Wish.AddWish,
		when.CallbackDataContains(c.CMD_ADD_TO_WISH),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveWish,
		when.MessageHasText(),
		with.AcceptFrom(addWish),
	)

	f.Activate(
		o.Wish.List,
		when.Command(c.CMD_WISHLIST),
	)

	f.Activate(
		o.Wish.List,
		when.CallbackDataContains(c.CMD_LIST),
	)

	f.Activate(
		o.Wish.WishInfo,
		when.CallbackDataContains(c.CMD_WISH_INFO),
	)

	f.Activate(
		o.Wish.DeleteWish,
		when.CallbackDataContains(c.CMD_DELETE_WISH),
	)

	f.Activate(
		o.Wish.ConfirmDeleteWish,
		when.CallbackDataContains(c.CMD_CONFIRM_DELETE_WISH),
	)

	editWishName := f.Activate(
		o.Wish.EditWishName,
		when.CallbackDataContains(c.CMD_EDIT_WISH_NAME),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveEditWishName,
		when.MessageHasText(),
		with.AcceptFrom(editWishName),
	)

	editWishLink := f.Activate(
		o.Wish.EditLink,
		when.CallbackDataContains(c.CMD_EDIT_LINK),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveEditLink,
		when.MessageHasText(),
		with.AcceptFrom(editWishLink),
	)

	editWishPrice := f.Activate(
		o.Wish.EditPrice,
		when.CallbackDataContains(c.CMD_EDIT_PRICE),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveEditPrice,
		when.MessageHasText(),
		with.AcceptFrom(editWishPrice),
	)

	f.Activate(
		o.Wish.DeleteLink,
		when.CallbackDataContains(c.CMD_DELETE_LINK),
		with.AcceptFrom(editWishLink),
	)

	f.Activate(
		o.Wish.ShareWishList,
		when.CallbackDataContains(c.CMD_SHARE_WISH_LIST),
	)

	f.Activate(
		o.Wish.ToggleWishLock,
		when.CallbackDataContains(c.CMD_TOGGLE_WISH_LOCK),
	)

	f.Activate(
		o.Wish.ShowSharedWishlist,
		appconditions.ShowSharedListCondition(),
	)

	f.Activate(
		o.Wish.WishInfo,
		when.CallbackDataContains(c.CMD_SHOW_SWI),
	)

	f.Activate(
		o.Wish.ShowSharedWishlist,
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
