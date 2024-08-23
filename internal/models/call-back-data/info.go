package models

import (
	"fmt"
	"strings"
)

type Info struct {
	id     string
	offset string
}

func CallInfo(id, offset string) *Info {
	return &Info{
		id: id,
		offset: offset,
	}
}

func InfoFromRaw(raw string) *Info {
	params := strings.Split(raw, ";")
	return &Info{
		id: params[1],
		offset: params[2],
	}
}

func (i *Info) String() string {
	return fmt.Sprintf("%s;%s;%s", i.Command(), i.ID(), i.Pagination().Offset)
}

func (i *Info) Command() string {
	return "info"
}

func (i *Info) ID() string {
	return i.id
}

func (i *Info) Entity() string {
	return "friend"
}

func (i *Info) Pagination() Pagination {
	return Pagination{Limit: "0", Offset: i.offset, Direction: ""}
}
