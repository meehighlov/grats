package inlinekeyboard

type button struct {
	text                            string
	callback_data                   string
	copy_text                       map[string]string
	switch_inline_query             string
	switch_inline_query_chosen_chat map[string]interface{}
	url                             string
	button_type                     string
}

func (b *Builder) NewButton(text, callback_data string) *button {
	return &button{text: text, callback_data: callback_data, button_type: "button"}
}

func (b *Builder) NewCopyButton(text, copy_text string) *button {
	return &button{
		text:        text,
		copy_text:   map[string]string{"text": copy_text},
		button_type: "copy",
	}
}

func (b *Builder) NewSwitchInlineButton(text string) *button {
	return &button{
		text:                text,
		switch_inline_query: "",
		button_type:         "switch_inline_query",
	}
}

func (b *Builder) NewAddToChatButton(text, query string) *button {
	return &button{
		text: text,
		switch_inline_query_chosen_chat: map[string]interface{}{
			"allow_user_chats":    false,
			"allow_bot_chats":     false,
			"allow_group_chats":   true,
			"allow_channel_chats": false,
			"query":               query,
		},
		button_type: "switch_inline_query_chosen_chat",
	}
}

func (b *Builder) NewURLButton(text, url string) *button {
	return &button{
		text:        text,
		url:         url,
		button_type: "url",
	}
}

func (b *Builder) NewShareLinkButton(text, link, description string) *button {
	query := link
	if description != "" {
		query = description + "\n\n" + link
	}
	return &button{
		text: text,
		switch_inline_query_chosen_chat: map[string]interface{}{
			"allow_user_chats":    true,
			"allow_bot_chats":     true,
			"allow_group_chats":   true,
			"allow_channel_chats": true,
			"query":               query,
		},
		button_type: "switch_inline_query_chosen_chat",
	}
}

func (b *button) Raw() *map[string]interface{} {
	result := map[string]interface{}{"text": b.text}

	if b.callback_data != "" {
		result["callback_data"] = b.callback_data
	}

	if b.copy_text != nil {
		result["copy_text"] = b.copy_text
	}

	if b.isInlineSwitchbutton() {
		result["switch_inline_query"] = b.switch_inline_query
	}

	if b.isChosenChatbutton() {
		result["switch_inline_query_chosen_chat"] = b.switch_inline_query_chosen_chat
	}

	if b.isURLbutton() {
		result["url"] = b.url
	}

	return &result
}

func (b *button) isInlineSwitchbutton() bool {
	return b.button_type == "switch_inline_query"
}

func (b *button) isChosenChatbutton() bool {
	return b.button_type == "switch_inline_query_chosen_chat"
}

func (b *button) isURLbutton() bool {
	return b.button_type == "url"
}

func (b *Builder) NewKeyboard() *Builder {
	return &Builder{
		keyboard: []*[]map[string]interface{}{},
	}
}

func (b *Builder) PutFirstAsLine(buttons ...*button) *Builder {
	rawbuttons := []map[string]interface{}{}
	for _, button := range buttons {
		rawbuttons = append(rawbuttons, *button.Raw())
	}

	b.keyboard = append([]*[]map[string]interface{}{&rawbuttons}, b.keyboard...)

	return b
}

// appends button list to representation of keyboard to new row below
func (b *Builder) AppendAsLine(buttons ...*button) *Builder {
	rawbuttons := []map[string]interface{}{}
	for _, button := range buttons {
		rawbuttons = append(rawbuttons, *button.Raw())
	}

	b.keyboard = append(b.keyboard, &rawbuttons)

	return b
}

// appends button list as stacked lines
func (b *Builder) AppendAsStack(buttons ...*button) *Builder {
	for _, button := range buttons {
		b.AppendAsLine(button)
	}

	return b
}

func (b *Builder) Murkup() []*[]map[string]interface{} {
	return b.keyboard
}

func (b *Builder) OptimizeView() *Builder {
	singleButtonRows := 0
	for _, row := range b.keyboard {
		if len(*row) == 1 {
			singleButtonRows++
		}
	}

	optimized := b.NewKeyboard()

	singleButtons := []*button{}

	for _, row := range b.keyboard {
		if len(*row) == 1 {
			rawButton := (*row)[0]
			btn := &button{
				text:        rawButton["text"].(string),
				button_type: "button",
			}
			if callbackData, ok := rawButton["callback_data"].(string); ok {
				btn.callback_data = callbackData
			}
			singleButtons = append(singleButtons, btn)
		} else {
			optimized.keyboard = append(optimized.keyboard, row)
		}
	}

	for i := 0; i < len(singleButtons); i += 2 {
		if i+1 < len(singleButtons) {
			optimized.AppendAsLine(singleButtons[i], singleButtons[i+1])
		} else {
			optimized.AppendAsLine(singleButtons[i])
		}
	}

	b.keyboard = optimized.keyboard

	return b
}

func (b *Builder) Append(keyboard *Builder) *Builder {
	b.keyboard = append(b.keyboard, keyboard.keyboard...)

	return b
}
