package common

import (
	"fmt"
	"strings"
)

type pagination struct {
	Offset    string
	Direction string
}

type CallbackDataModel struct {
	Command    string
	Id         string
	Pagination pagination
	Entity     string
	SourceId   string
}

type ListCaller func(offset, direction, sourceId, entity string) *CallbackDataModel

func CallList(offset, direction, sourceId string, entity string) *CallbackDataModel {
	return newCallback("list", "", offset, direction, entity, sourceId)
}

func CallDelete(id, offset string) *CallbackDataModel {
	return newCallback("delete", id, offset, "", "friend", "")
}

func CallInfo(id, offset, table string) *CallbackDataModel {
	return newCallback(fmt.Sprintf("%s_info", table), id, offset, "", table, "")
}

func CallChatInfo(id string) *CallbackDataModel {
	return newCallback("chat_info", id, "", "", "chat", "")
}

func CallChatList() *CallbackDataModel {
	return newCallback("chat_list", "", "", "", "chat", "")
}

func CallAddItem(id, table string) *CallbackDataModel {
	return newCallback(fmt.Sprintf("add_to_%s", table), id, "", "", table, "")
}

func CallChatHowto(id string) *CallbackDataModel {
	return newCallback("chat_howto", id, "", "", "chat", "")
}

func CallEditGreetingTemplate(id string) *CallbackDataModel {
	return newCallback("edit_greeting_template", id, "", "", "chat", "")
}

func CallDeleteChat(id string) *CallbackDataModel {
	return newCallback("delete_chat", id, "", "", "chat", "")
}

func CallConfirmDeleteChat(id string) *CallbackDataModel {
	return newCallback("confirm_delete_chat", id, "", "", "chat", "")
}

func CallToggleSilentNotifications(id string) *CallbackDataModel {
	return newCallback("toggle_silent_notifications", id, "", "", "chat", "")
}

func CallShareWishList(chatId string) *CallbackDataModel {
	return newCallback("share_wish_list", chatId, "", "", "wish", "")
}

func CallDeleteWish(id, offset string) *CallbackDataModel {
	return newCallback("delete_wish", id, offset, "", "wish", "")
}

func CallEditPrice(id string) *CallbackDataModel {
	return newCallback("edit_price", id, "", "", "wish", "")
}

func CallEditLink(id string) *CallbackDataModel {
	return newCallback("edit_link", id, "", "", "wish", "")
}

func CallWishInfo(id, offset string) *CallbackDataModel {
	return newCallback("wish_info", id, offset, "", "wish", "")
}

func CallSharedWishInfo(id, offset string) *CallbackDataModel {
	return newCallback("show_swi", id, offset, "", "wish", "")
}

func CallSharedWishList(offset, direction, sourceId, entity string) *CallbackDataModel {
	return newCallback("show_swl", sourceId, offset, direction, entity, "")
}

func CallToggleWishLock(id, offset string) *CallbackDataModel {
	return newCallback("toggle_wish_lock", id, offset, "", "wish", "")
}

func CallConfirmDeleteWish(id string) *CallbackDataModel {
	return newCallback("confirm_delete_wish", id, "", "", "wish", "")
}

func CallCommands() *CallbackDataModel {
	return newCallback("commands", "", "", "", "", "")
}

func CallConfirmDelete(id string) *CallbackDataModel {
	return newCallback("confirm_delete", id, "", "", "friend", "")
}

func CallEditName(id string) *CallbackDataModel {
	return newCallback("edit_name", id, "", "", "friend", "")
}

func CallEditBirthday(id string) *CallbackDataModel {
	return newCallback("edit_birthday", id, "", "", "friend", "")
}

func CallEditWishName(id string) *CallbackDataModel {
	return newCallback("edit_wish_name", id, "", "", "wish", "")
}

func CallSupport(chatId string) *CallbackDataModel {
	return newCallback("support", chatId, "", "", "support", "")
}

func CallWriteToSupport(chatId string) *CallbackDataModel {
	return newCallback("write_to_support", chatId, "", "", "support", "")
}

func newCallback(command, id, offset, direction, entity, sourceId string) *CallbackDataModel {
	return &CallbackDataModel{
		Command: command,
		Id:      id,
		Pagination: pagination{
			Offset:    offset,
			Direction: direction,
		},
		Entity:   entity,
		SourceId: sourceId,
	}
}

func CallbackFromString(raw string) *CallbackDataModel {
	params := strings.Split(raw, ";")
	sourceId := ""
	if len(params) == 6 {
		sourceId = params[5]
	}
	return &CallbackDataModel{
		Command: params[0],
		Id:      params[1],
		Pagination: pagination{
			Offset:    params[2],
			Direction: params[3],
		},
		Entity:   params[4],
		SourceId: sourceId,
	}
}

func (cd *CallbackDataModel) String() string {
	separator := ";"
	return strings.Join(
		[]string{
			cd.Command,
			cd.Id,
			cd.Pagination.Offset,
			cd.Pagination.Direction,
			cd.Entity,
			cd.SourceId,
		},
		separator,
	)
}
