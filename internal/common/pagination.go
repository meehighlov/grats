package common

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"gorm.io/gorm"
)

type SearchParams struct {
	ListId string
}

type PaginatedEntity interface {
	GetId() string
	GetUserId() string
	GreaterThan(other PaginatedEntity) bool
	ButtonText() string
	Search(ctx context.Context, tx *gorm.DB, params *SearchParams) ([]PaginatedEntity, error)
}

func BuildEntityButtons[T PaginatedEntity](
	entities []T,
	limit,
	offset int,
	callbackDataBuilder func(id string, offset int) string,
) []*Button {
	sort.Slice(entities, func(i, j int) bool { return comparator(entities, i, j) })
	buttons := []*Button{}
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

		buttons = append(buttons, NewButton(buttonText, callbackDataBuilder(entity.GetId(), offset)))
	}

	return buttons
}

func GetOffsetByDirection(
	direction string,
	offset int,
	entities []PaginatedEntity,
	paginationShift int,
) int {
	offset_ := offset
	if direction == "<" {
		return offset_
	}
	if direction == "<<<" {
		offset_ = 0
	}
	if direction == ">>" {
		offset_ += paginationShift
	}
	if direction == "<<" {
		offset_ -= paginationShift
	}
	if direction == "<>" {
		offset_ = len(entities)
	}
	return offset_
}

func comparator[T PaginatedEntity](entities []T, i, j int) bool {
	return (entities[i]).GreaterThan(entities[j])
}

func AppendControlButtons(
	keyboard *InlineKeyboard,
	total,
	limit,
	offset int,
	listId string,
	entity string,
	listCaller ListCaller,
	listStartOffset int,
) error {
	buttons := []*Button{}
	if total <= limit || total == 0 {
		return nil
	}
	if offset == total {
		buttons = append(buttons, NewButton("⬆️", listCaller(strconv.Itoa(listStartOffset), "<<<", listId, entity).String()))
		keyboard.AppendAsLine(buttons...)
		return nil
	}
	if offset+limit >= total {
		buttons = append(buttons, NewButton("⬅️", listCaller(strconv.Itoa(offset), "<<", listId, entity).String()))
	} else {
		if offset == 0 {
			buttons = append(buttons, NewButton("➡️", listCaller(strconv.Itoa(offset), ">>", listId, entity).String()))
		} else {
			buttons = append(buttons, NewButton("⬅️", listCaller(strconv.Itoa(offset), "<<", listId, entity).String()))
			buttons = append(buttons, NewButton("➡️", listCaller(strconv.Itoa(offset), ">>", listId, entity).String()))
		}
	}

	keyboard.AppendAsLine(buttons...)
	keyboard.AppendAsStack(NewButton(fmt.Sprintf("(%d)⬇️", total), listCaller(strconv.Itoa(offset), "<>", listId, entity).String()))

	return nil
}
