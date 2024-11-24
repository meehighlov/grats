package common

type Button struct {
	text          string
	callback_data string
}

func NewButton(text, callback_data string) *Button {
	return &Button{text, callback_data}
}

func (b *Button) Raw() map[string]string {
	return map[string]string{"text": b.text, "callback_data": b.callback_data}
}

type InlineKeyboard struct {
	markup [][]map[string]string
}

func NewInlineKeyboard() *InlineKeyboard {
	return &InlineKeyboard{[][]map[string]string{}}
}

// appends button list to representation of keyboard to new row below
func (ik *InlineKeyboard) AppendAsLine(buttons ...Button) {
	rawButtons := []map[string]string{}
	for _, button := range buttons {
		rawButtons = append(rawButtons, button.Raw())
	}

	ik.markup = append(ik.markup, rawButtons)
}

// appends button list as stacked lines
func (ik *InlineKeyboard) AppendAsStack(buttons ...Button) {
	for _, button := range buttons {
		ik.AppendAsLine(button)
	}
}

func (ik *InlineKeyboard) Murkup() *[][]map[string]string {
	return &ik.markup
}
