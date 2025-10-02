package callbackdata

import "strings"

type callbackData struct {
	ID      string
	Command string
	Offset  string
}

func (b *Builder) Build(id, command, offset string) *callbackData {
	return &callbackData{
		ID:      id,
		Command: command,
		Offset:  offset,
	}
}

func (b *Builder) FromString(raw string) *callbackData {
	params := strings.Split(raw, ";")
	if len(params) != 3 {
		return &callbackData{
			Offset: "0",
		}
	}

	return &callbackData{
		Command: params[0],
		ID:      params[1],
		Offset:  params[2],
	}
}

func (cd *callbackData) String() string {
	separator := ";"
	return strings.Join(
		[]string{
			cd.Command,
			cd.ID,
			cd.Offset,
		},
		separator,
	)
}
