package models

type CallbackDataModel interface {
	String() string  // string cast
	Command() string  // handler to call in callback query handler
	Pagination() (limit, offset, direction string)  // pagination
	ID() // id of enitity
	Entity() string  // entity name which id belongs to
}
