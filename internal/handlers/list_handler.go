package handlers

import (
	"context"
	"fmt"
	"strconv"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
)

const (
	LIST_PAGINATION_SHIFT = 5
	LIST_LIMIT            = 5
	LIST_START_OFFSET     = 0

	HEADER_MESSAGE_LIST_NOT_EMPTY = "✨Список %s"
	HEADER_MESSAGE_LIST_IS_EMPTY  = "✨Список пуст"
)

func ListItemsHandler(ctx context.Context, event *common.Event) error {
	callbackData := common.CallbackFromString(event.GetCallbackQuery().Data)

	listId := callbackData.ListId
	entity := callbackData.Entity

	entities, err := db.CreateEntity(entity).Search(ctx, nil, &common.SearchParams{
		ListId: listId,
	})

	if err != nil {
		event.Logger.Error("Error fetching items" + err.Error())
		return err
	}

	offset, _ := strconv.Atoi(callbackData.Pagination.Offset)
	if offset == 0 {
		offset = LIST_START_OFFSET
	}
	direction := callbackData.Pagination.Direction

	offset_ := common.GetOffsetByDirection(direction, offset, entities, LIST_PAGINATION_SHIFT)

	if _, err := event.EditCalbackMessage(
		ctx,
		buildChatHeaderMessage(ctx, listId, event, entities, entity),
		*buildListMarkup(entities, LIST_LIMIT, offset_, listId, entity).Murkup(),
	); err != nil {
		return err
	}

	return nil
}

func buildListMarkup(entities []common.PaginatedEntity, limit, offset int, listId string, table string) *common.InlineKeyboard {
	callbackDataBuilder := func(id string, offset int) string {
		return common.CallInfo(id, strconv.Itoa(offset), table).String()
	}
	entityListAsButtons := common.BuildEntityButtons(entities, limit, offset, callbackDataBuilder)
	keyboard := common.NewInlineKeyboard()

	headerButtons := []*common.Button{
		common.NewButton("↩️", common.CallCommands().String()),
		common.NewButton("➕", common.CallAddItem(listId, table).String()),
	}

	if table == "wish" && len(entities) > 0 {
		headerButtons = append(headerButtons, common.NewButton("🛜", common.CallShareWishList(listId).String()))
	}

	keyboard.AppendAsLine(headerButtons...)
	keyboard.AppendAsStack(entityListAsButtons...)

	common.AppendControlButtons(keyboard, len(entities), limit, offset, listId, table, common.CallList, LIST_START_OFFSET)

	return keyboard
}

func buildChatHeaderMessage(ctx context.Context, chatId string, event *common.Event, entities []common.PaginatedEntity, table string) string {
	if len(entities) == 0 {
		return HEADER_MESSAGE_LIST_IS_EMPTY
	}
	if table == "wish" {
		userId := entities[0].GetUserId()
		userInfo, _ := event.GetChatMember(ctx, userId)
		return fmt.Sprintf("✨Вишлист @%s", userInfo.Result.User.Username)
	}

	return HEADER_MESSAGE_LIST_IS_EMPTY
}
