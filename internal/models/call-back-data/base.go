package models

import "strings"

type Pagination struct {
	Limit     string
	Offset    string
	Direction string
}

type CallbackDataModel interface {
	String() string  // string cast
	Command() string  // handler to call in callback query handler
	Pagination() Pagination  // pagination
	ID() // id of enitity
	Entity() string  // entity name which id belongs to
}

func GetCommandFromCallbackData(callbackData string) string {
	return strings.Split(callbackData, ";")[0]
}
