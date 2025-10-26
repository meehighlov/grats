package pagination

import (
	"strconv"

	inlinekeyboard "github.com/meehighlov/grats/pkg/telegram/builders/inline_keyboard"
)

func (p *Builder) BuildControls(total int, command, chat_id string, offset int) *inlinekeyboard.Builder {
	limit := p.Limit
	baseOffset := p.BaseOffset
	callback_data_builder := p.callbackDataBuilder
	keyboard := p.keyboardBuilder.NewKeyboard()

	if total <= limit {
		return keyboard
	}

	if offset == 0 {
		callback_data := callback_data_builder.Build(chat_id, command, strconv.Itoa(baseOffset))
		return keyboard.AppendAsLine(keyboard.NewButton("➡️", callback_data.String()))
	}

	if offset+limit >= total {
		callback_data := callback_data_builder.Build(chat_id, command, strconv.Itoa(offset-baseOffset))
		return keyboard.AppendAsLine(keyboard.NewButton("⬅️", callback_data.String()))
	}

	if offset > 0 {
		callback_data_prev := callback_data_builder.Build(chat_id, command, strconv.Itoa(offset-baseOffset))
		callback_data_next := callback_data_builder.Build(chat_id, command, strconv.Itoa(offset+baseOffset))

		keyboard.AppendAsLine(
			keyboard.NewButton("⬅️", callback_data_prev.String()),
			keyboard.NewButton("➡️", callback_data_next.String()),
		)
	}

	return keyboard
}
