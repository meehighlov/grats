package handlers

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/db"
	"gorm.io/gorm"
)

const (
	LIST_PAGINATION_SHIFT = 5
	LIST_LIMIT            = 5
	LIST_START_OFFSET     = 0

	HEADER_MESSAGE_LIST_NOT_EMPTY = "✨Список %s"
	HEADER_MESSAGE_LIST_IS_EMPTY  = "✨Записей пока нет"
)

func ListItemsHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	callbackData := common.CallbackFromString(event.GetCallbackQuery().Data)

	chatId := callbackData.Id
	userId := strconv.Itoa(event.GetCallbackQuery().From.Id)
	entity := callbackData.Entity

	entities, err := db.NewEntity(entity).Search(ctx, tx, chatId, userId)

	if err != nil {
		event.Logger.Error("Error fetching friends" + err.Error())
		return err
	}

	if event.GetCallbackQuery().Id != "" {
		if _, err := event.EditCalbackMessage(
			ctx,
			buildChatHeaderMessage(ctx, chatId, event, (len(entities) == 0), entity),
			*buildListMarkup(entities, LIST_LIMIT, LIST_START_OFFSET, chatId, entity).Murkup(),
		); err != nil {
			return err
		}

		return nil
	}

	if _, err := event.ReplyWithKeyboard(
		ctx,
		buildChatHeaderMessage(ctx, chatId, event, (len(entities) == 0), entity),
		*buildListMarkup(entities, LIST_LIMIT, LIST_START_OFFSET, chatId, entity).Murkup(),
	); err != nil {
		return err
	}

	return nil
}

func comparator[T db.Entity](entities []T, i, j int) bool {
	return (entities[i]).GreaterThan(entities[j])
}

func appendControlButtons(keyboard *common.InlineKeyboard, total, limit, offset int, chatId string) error {
	buttons := []*common.Button{}
	if total <= limit || total == 0 {
		return nil
	}
	if offset == total {
		buttons = append(buttons, common.NewButton("⬆️", common.CallList(strconv.Itoa(LIST_START_OFFSET), "<<<", chatId).String()))
		keyboard.AppendAsLine(buttons...)
		return nil
	}
	if offset+limit >= total {
		buttons = append(buttons, common.NewButton("⬅️", common.CallList(strconv.Itoa(offset), "<<", chatId).String()))
	} else {
		if offset == 0 {
			buttons = append(buttons, common.NewButton("➡️", common.CallList(strconv.Itoa(offset), ">>", chatId).String()))
		} else {
			buttons = append(buttons, common.NewButton("⬅️", common.CallList(strconv.Itoa(offset), "<<", chatId).String()))
			buttons = append(buttons, common.NewButton("➡️", common.CallList(strconv.Itoa(offset), ">>", chatId).String()))
		}
	}

	keyboard.AppendAsLine(buttons...)
	keyboard.AppendAsStack(common.NewButton(fmt.Sprintf("(%d)⬇️", total), common.CallList(strconv.Itoa(offset), "<>", chatId).String()))

	return nil
}

func ListPaginationCallbackQueryHandler(ctx context.Context, event *common.Event, tx *gorm.DB) error {
	params := common.CallbackFromString(event.GetCallbackQuery().Data)

	offset := params.Pagination.Offset

	limit_ := LIST_LIMIT
	offset_, err := strconv.Atoi(offset)
	if err != nil {
		event.Logger.Error("error parsing offset in list pagination callback query: " + err.Error())
		return err
	}

	chatId := params.BoundChat

	entities, err := db.NewEntity(params.Entity).Search(ctx, tx, chatId, "")

	if err != nil {
		event.Logger.Error("Error fetching friends" + err.Error())
		return err
	}

	direction := params.Pagination.Direction

	event.Logger.Debug(fmt.Sprintf("direction: %s limit: %d offset: %s", direction, limit_, offset))

	if direction == "<" {
		event.Logger.Debug("back to previous screen, offset not changed")
	}
	if direction == "<<<" {
		offset_ = 0
	}
	if direction == ">>" {
		offset_ += LIST_PAGINATION_SHIFT
	}
	if direction == "<<" {
		offset_ -= LIST_PAGINATION_SHIFT
	}
	if direction == "<>" {
		offset_ = len(entities)
	}

	msg := buildChatHeaderMessage(ctx, chatId, event, (len(entities) == 0), params.Entity)

	if _, err := event.EditCalbackMessage(ctx, msg, *buildListMarkup(entities, limit_, offset_, chatId, params.Entity).Murkup()); err != nil {
		return err
	}

	return nil
}

func buildEntityButtons[T db.Entity](entities []T, limit, offset int, callbackDataBuilder func(id string, offset int) string) []*common.Button {
	sort.Slice(entities, func(i, j int) bool { return comparator(entities, i, j) })
	buttons := []*common.Button{}
	for i, entity := range entities {
		if offset != len(entities) {
			if i == limit+offset {
				break
			}
			if i < offset {
				continue
			}
		}

		buttonText := entity.ButtonText()

		buttons = append(buttons, common.NewButton(buttonText, callbackDataBuilder(entity.GetId(), offset)))
	}

	return buttons
}

func buildListMarkup[T db.Entity](entities []T, limit, offset int, chatId string, table string) *common.InlineKeyboard {
	callbackDataBuilder := func(id string, offset int) string {
		return common.CallInfo(id, strconv.Itoa(offset), table).String()
	}
	entityListAsButtons := buildEntityButtons(entities, limit, offset, callbackDataBuilder)
	keyboard := common.NewInlineKeyboard()

	keyboard.AppendAsLine(common.NewButton("🏠 в начало", common.CallSetup().String()))
	keyboard.AppendAsLine(common.NewButton("➕ добавить", common.CallAddItem(chatId, table).String()))

	keyboard.AppendAsStack(entityListAsButtons...)

	appendControlButtons(keyboard, len(entities), limit, offset, chatId)

	if strings.Contains(chatId, "-") {
		keyboard.AppendAsLine(common.NewButton("⬅️к чату", common.CallChatInfo(chatId).String()))
	}

	return keyboard
}

func buildChatHeaderMessage(ctx context.Context, chatId string, event *common.Event, emptyList bool, table string) string {
	if emptyList {
		return HEADER_MESSAGE_LIST_IS_EMPTY
	}
	if table == "wish" {
		return fmt.Sprintf(HEADER_MESSAGE_LIST_NOT_EMPTY, "желаний")
	}
	chatFullInfo, err := event.GetChat(ctx, chatId)
	if err != nil {
		return fmt.Sprintf(HEADER_MESSAGE_LIST_NOT_EMPTY, "напоминаний о др в чате")
	}
	if chatFullInfo.Id < 0 {
		return fmt.Sprintf(HEADER_MESSAGE_LIST_NOT_EMPTY, "напоминаний о др в чате "+chatFullInfo.Title)
	} else {
		return fmt.Sprintf(HEADER_MESSAGE_LIST_NOT_EMPTY, "личных напоминаний о др")
	}
}
