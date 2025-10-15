package app

import (
	"context"

	appconditions "github.com/meehighlov/grats/internal/app/conditions"
	"github.com/meehighlov/grats/internal/clients"
	"github.com/meehighlov/grats/internal/clients/clients/telegram"
	"github.com/meehighlov/grats/internal/config"
	"github.com/meehighlov/grats/internal/constants"
	"github.com/meehighlov/grats/internal/errs"
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

	f.Activate(
		o.Support.SupportHandler,
		when.Command(c.CMD_SUPPORT),
	)

	f.Activate(
		o.Support.WriteHandler,
		when.CallbackDataContains(c.CMD_SUPPORT_WRITE),
		with.Transition(nil, c.CMD_SUPPORT_SEND),
		with.BeforeHandler(resetUserCache),
	)

	f.Activate(
		o.Support.SendMessageHandler,
		when.UpdateHasOnlyText(),
		with.AllowedActivationStatus(c.CMD_SUPPORT_SEND),
	)

	f.Activate(
		o.Support.CancelHandler,
		when.CallbackDataContains(c.CMD_SUPPORT_CANCEL),
	)

	f.Activate(
		o.Support.HandleSupportReply,
		appconditions.SupportReplyCondition(cfg.SupportChatId),
	)

	f.Activate(
		o.Wish.AddWishHandler,
		when.CallbackDataContains(c.CMD_ADD_TO_WISH),
		with.Transition(nil, c.CMD_ADD_SAVE_WISH),
		with.BeforeHandler(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveWish,
		when.UpdateHasOnlyText(),
		with.AllowedActivationStatus(c.CMD_ADD_SAVE_WISH),
		with.Transition(errs.ErrSaveWishValidation, c.CMD_ADD_SAVE_WISH),
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
		o.Wish.WishInfoHandler,
		when.CallbackDataContains(c.CMD_WISH_INFO),
	)

	f.Activate(
		o.Wish.DeleteWishCallbackQueryHandler,
		when.CallbackDataContains(c.CMD_DELETE_WISH),
	)

	f.Activate(
		o.Wish.ConfirmDeleteWishCallbackQueryHandler,
		when.CallbackDataContains(c.CMD_CONFIRM_DELETE_WISH),
	)

	f.Activate(
		o.Wish.EditWishNameHandler,
		when.CallbackDataContains(c.CMD_EDIT_WISH_NAME),
		with.Transition(nil, c.CMD_EDIT_WISH_NAME_SAVE),
		with.BeforeHandler(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveEditWishNameHandler,
		when.UpdateHasOnlyText(),
		with.AllowedActivationStatus(c.CMD_EDIT_WISH_NAME_SAVE),
		with.Transition(errs.ErrEditWishNameValidation, c.CMD_EDIT_WISH_NAME_SAVE),
	)

	f.Activate(
		o.Wish.EditLinkHandler,
		when.CallbackDataContains(c.CMD_EDIT_LINK),
		with.Transition(nil, c.CMD_EDIT_LINK_SAVE),
		with.BeforeHandler(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveEditLinkHandler,
		when.UpdateHasOnlyText(),
		with.AllowedActivationStatus(c.CMD_EDIT_LINK_SAVE),
		with.Transition(errs.ErrSaveEditLinkValidation, c.CMD_EDIT_LINK_SAVE),
	)

	f.Activate(
		o.Wish.EditPriceHandler,
		when.CallbackDataContains(c.CMD_EDIT_PRICE),
		with.Transition(nil, c.CMD_EDIT_PRICE_SAVE),
		with.BeforeHandler(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveEditPriceHandler,
		when.UpdateHasOnlyText(),
		with.AllowedActivationStatus(c.CMD_EDIT_PRICE_SAVE),
		with.Transition(errs.ErrEditPriceValidation, c.CMD_EDIT_PRICE_SAVE),
	)

	f.Activate(
		o.Wish.DeleteLinkHandler,
		when.CallbackDataContains(c.CMD_DELETE_LINK),
	)

	f.Activate(
		o.Wish.ShareWishListHandler,
		when.CallbackDataContains(c.CMD_SHARE_WISH_LIST),
	)

	f.Activate(
		o.Wish.ToggleWishLockHandler,
		when.CallbackDataContains(c.CMD_TOGGLE_WISH_LOCK),
	)

	f.Activate(
		o.Wish.ShowSharedWishlistHandler,
		appconditions.ShowSharedListCondition(),
	)

	f.Activate(
		o.Wish.WishInfoHandler,
		when.CallbackDataContains(c.CMD_SHOW_SWI),
	)

	f.Activate(
		o.Wish.ShowSharedWishlistHandler,
		when.CallbackDataContains(c.CMD_SHOW_SWL),
	)
}
