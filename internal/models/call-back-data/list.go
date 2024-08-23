package models

import (
	"fmt"
	"strings"
)

type List struct {
	offset     string
	direction  string
}

func CallList(offset, direction string) *List {
	return &List{
		offset: offset,
		direction: direction,
	}
}

func ListFromRaw(raw string) *List {
	params := strings.Split(raw, ";")
	return &List{
		offset: params[1],
		direction: params[2],
	}
}

func (l *List) String() string {
	return fmt.Sprintf("%s;%s;%s", l.Command(), l.offset, l.direction)
}

func (*List) Command() string {
	return "list"
}

func (l *List) Pagination() Pagination {
	return Pagination{Limit: "0", Offset: l.offset, Direction: l.direction}
}

func (*List) ID() string {
	return ""
}
