package models

import (
	"fmt"
	"strings"
)

type Delete struct {
	id     string
}

func СallDelete(id string) *Delete {
	return &Delete{
		id: id,
	}
}

func DeleteFromRaw(raw string) *Delete {
	params := strings.Split(raw, ";")
	return &Delete{
		id: params[1],
	}
}

func (d *Delete) String() string {
	return fmt.Sprintf("%s;%s", d.Command(), d.id)
}

func (*Delete) Command() string {
	return "delete"
}

func (*Delete) Entity() string {
	return "friend"
}

func (d *Delete) Pagination() Pagination {
	return Pagination{Limit: "0", Offset: "0", Direction: ""}
}

func (d *Delete) ID() string {
	return d.id
}
