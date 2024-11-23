package common

import (
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
	BoundChat  string
}

func CallList(offset, direction, chatid string) *CallbackDataModel {
	return newCallback("list", "", offset, direction, "friend", chatid)
}

func CallDelete(id, chatid string) *CallbackDataModel {
	return newCallback("delete", id, "", "", "friend", chatid)
}

func CallInfo(id, offset, chatid string) *CallbackDataModel {
	return newCallback("info", id, offset, "", "friend", chatid)
}

func CallChatInfo(id string) *CallbackDataModel {
	return newCallback("chat_info", id, "", "", "chat", id)
}

func CallChatList() *CallbackDataModel {
	return newCallback("chat_list", "", "", "", "chat", "")
}

func CallAddToChat(id string) *CallbackDataModel {
	return newCallback("add_to_chat", id, "", "", "chat", id)
}

func CallChatBirthdays(id string) *CallbackDataModel {
	return newCallback("chat_birthdays", id, "", "", "chat", id)
}

func CallChatHowto(id string) *CallbackDataModel {
	return newCallback("chat_howto", id, "", "", "chat", id)
}

func newCallback(command, id, offset, direction, entity, chatid string) *CallbackDataModel {
	return &CallbackDataModel{
		Command: command,
		Id:      id,
		Pagination: pagination{
			Offset:    offset,
			Direction: direction,
		},
		Entity:    entity,
		BoundChat: chatid,
	}
}

func CallbackFromString(raw string) *CallbackDataModel {
	params := strings.Split(raw, ";")
	boundChat := ""
	if len(params) == 6 {
		boundChat = params[5]
	}
	return &CallbackDataModel{
		Command: params[0],
		Id:      params[1],
		Pagination: pagination{
			Offset:    params[2],
			Direction: params[3],
		},
		Entity:    params[4],
		BoundChat: boundChat,
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
			cd.BoundChat,
		},
		separator,
	)
}
