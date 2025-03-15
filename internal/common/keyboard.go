package common

type Button struct {
	text          string
	callback_data string
	copy_text     map[string]string
}

func NewButton(text, callback_data string) *Button {
	return &Button{text: text, callback_data: callback_data}
}

func NewCopyButton(text, copy_text string) *Button {
	return &Button{
		text:      text,
		copy_text: map[string]string{"text": copy_text},
	}
}

func (b *Button) Raw() map[string]interface{} {
	result := map[string]interface{}{"text": b.text}

	if b.callback_data != "" {
		result["callback_data"] = b.callback_data
	}

	if b.copy_text != nil {
		result["copy_text"] = b.copy_text
	}

	return result
}

type InlineKeyboard struct {
	markup [][]map[string]interface{}
}

func NewInlineKeyboard() *InlineKeyboard {
	return &InlineKeyboard{[][]map[string]interface{}{}}
}

// appends button list to representation of keyboard to new row below
func (ik *InlineKeyboard) AppendAsLine(buttons ...Button) {
	rawButtons := []map[string]interface{}{}
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

func (ik *InlineKeyboard) Murkup() *[][]map[string]interface{} {
	return &ik.markup
}
