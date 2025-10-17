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
		o.Support.Support,
		when.Command(c.CMD_SUPPORT),
	)

	f.Activate(
		o.Support.SupportWrite,
		when.CallbackDataContains(c.CMD_SUPPORT_WRITE),
		with.ID(c.CMD_SUPPORT_WRITE),
		with.SuccessOutput(c.CMD_SUPPORT_SEND),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Support.SendSupportMessage,
		when.UpdateHasOnlyText(),
		with.ID(c.CMD_SUPPORT_SEND),
		with.InputState(c.CMD_SUPPORT_WRITE),
	)

	f.Activate(
		o.Support.CancelSupportCall,
		when.CallbackDataContains(c.CMD_SUPPORT_CANCEL),
	)

	f.Activate(
		o.Support.ProcessSupportReply,
		appconditions.SupportReplyCondition(cfg.SupportChatId),
	)

	f.Activate(
		o.Wish.AddWish,
		when.CallbackDataContains(c.CMD_ADD_TO_WISH),
		with.ID(c.CMD_ADD_TO_WISH),
		with.SuccessOutput(c.CMD_ADD_SAVE_WISH),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveWish,
		when.UpdateHasOnlyText(),
		with.ID(c.CMD_ADD_SAVE_WISH),
		with.InputState(c.CMD_ADD_TO_WISH),
		with.Retry(errs.ErrSaveWishValidation),
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

	f.Activate(
		o.Wish.EditWishName,
		when.CallbackDataContains(c.CMD_EDIT_WISH_NAME),
		with.ID(c.CMD_EDIT_WISH_NAME),
		with.SuccessOutput(c.CMD_EDIT_WISH_NAME_SAVE),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveEditWishName,
		when.UpdateHasOnlyText(),
		with.ID(c.CMD_EDIT_WISH_NAME_SAVE),
		with.InputState(c.CMD_EDIT_WISH_NAME),
		with.Retry(errs.ErrEditWishNameValidation),
	)

	f.Activate(
		o.Wish.EditLink,
		when.CallbackDataContains(c.CMD_EDIT_LINK),
		with.ID(c.CMD_EDIT_LINK),
		with.SuccessOutput(c.CMD_EDIT_LINK_SAVE),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveEditLink,
		when.UpdateHasOnlyText(),
		with.ID(c.CMD_EDIT_LINK_SAVE),
		with.InputState(c.CMD_EDIT_LINK),
		with.Retry(errs.ErrSaveEditLinkValidation),
	)

	f.Activate(
		o.Wish.EditPrice,
		when.CallbackDataContains(c.CMD_EDIT_PRICE),
		with.ID(c.CMD_EDIT_PRICE),
		with.SuccessOutput(c.CMD_EDIT_PRICE_SAVE),
		with.BeforeAction(resetUserCache),
	)

	f.Activate(
		o.Wish.SaveEditPrice,
		when.UpdateHasOnlyText(),
		with.ID(c.CMD_EDIT_PRICE_SAVE),
		with.InputState(c.CMD_EDIT_PRICE),
		with.Retry(errs.ErrEditPriceValidation),
	)

	f.Activate(
		o.Wish.DeleteLink,
		when.CallbackDataContains(c.CMD_DELETE_LINK),
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
}
