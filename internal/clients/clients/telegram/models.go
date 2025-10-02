package telegram

import (
	"os"
	"reflect"
	"strconv"
	"strings"
)

type Chat struct {
	// full description https://core.telegram.org/bots/api#chat
	//Type of chat, can be either "private", "group", "supergroup" or "channel"
	Id        int    `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
}

type User struct {
	// full description https://core.telegram.org/bots/api#user
	Id        int    `json:"id"`
	IsBot     bool   `json:"is_bot"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
}

func (user *User) IsAdmin() bool {
	for _, admin_name := range strings.Split(os.Getenv("ADMINS"), ",") {
		if user.Username == admin_name {
			return true
		}
	}

	return false
}

type ChatMember struct {
	// full description https://core.telegram.org/bots/api#chatmemberowner
	Status string `json:"status"`
	User   User   `json:"user"`
}

type ChatMemberResponse struct {
	Ok     bool         `json:"ok"`
	Result []ChatMember `json:"result"`
}

type SingleChatMemberResponse struct {
	Ok     bool       `json:"ok"`
	Result ChatMember `json:"result"`
}

type ReplyToMessage struct {
	MessageId  int    `json:"message_id"`
	From       User   `json:"from"`
	SenderChat Chat   `json:"sender_chat"`
	Chat       Chat   `json:"chat"`
	Text       string `json:"text"`
}

type Message struct {
	MessageId      int                  `json:"message_id"`
	From           User                 `json:"from"`
	SenderChat     Chat                 `json:"sender_chat"`
	Chat           Chat                 `json:"chat"`
	Text           string               `json:"text"`
	ReplyToMessage ReplyToMessage       `json:"reply_to_message"`
	NewChatMembers []User               `json:"new_chat_members"`
	LeftChatMember User                 `json:"left_chat_member"`
	ReplyMarkup    InlineKeyboardMarkup `json:"reply_markup"`
}

func (m *Message) IsReply() bool {
	return reflect.ValueOf(m).Elem().FieldByName("ReplyToMessage") != reflect.Value{}
}

func (m *Message) HasLeftChatMember() bool {
	return reflect.ValueOf(m).Elem().FieldByName("LeftChatMember") != reflect.Value{}
}

type Update struct {
	UpdateId      int           `json:"update_id"`
	Message       Message       `json:"message"`
	InlineQuery   InlineQuery   `json:"inline_query"`
	CallbackQuery CallbackQuery `json:"callback_query"`
}

type InlineQuery struct {
	Id    string `json:"id"`
	From  User   `json:"from"`
	Query string `json:"query"`
}

type CallbackQuery struct {
	Id      string  `json:"id"`
	From    User    `json:"from"`
	Message Message `json:"message"`
	Data    string  `json:"data"`
}

type GetMeReponse struct {
	Ok     bool `json:"ok"`
	Result User `json:"result"`
}

type UpdateResponse struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type GetChatResponse struct {
	Ok     bool `json:"ok"`
	Result Chat `json:"result"`
}

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	CallbackData string `json:"callback_data"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboard [][]InlineKeyboardButton `json:"inline_keyboard"`
}

func (update *UpdateResponse) GetLastUpdateId() int {
	return update.Result[len(update.Result)-1].UpdateId
}

func (message *Message) GetChatIdStr() string {
	return strconv.Itoa(message.Chat.Id)
}

func (message *Message) GetSenderChatIdStr() string {
	return strconv.Itoa(message.SenderChat.Id)
}

func (message *Message) GetMessageIdStr() string {
	return strconv.Itoa(message.MessageId)
}

func (message *Message) GetCommand() string {
	if strings.HasPrefix(message.Text, "/") {
		return message.Text
	}
	return ""
}

func (u *Update) GetChatIdStr() string {
	return u.GetMessage().GetChatIdStr()
}

func (u *Update) IsCallback() bool {
	return u.CallbackQuery.Id != ""
}

func (u *Update) GetMessageIdStr() string {
	return u.GetMessage().GetMessageIdStr()
}

func (u *Update) GetMessage() *Message {
	if u.Message.MessageId != 0 {
		return &u.Message
	}

	if u.IsCallback() {
		return &u.CallbackQuery.Message
	}

	return &u.Message
}

type Document struct {
	FileId       string `json:"file_id"`
	FileUniqueId string `json:"file_unique_id"`
	FileName     string `json:"file_name"`
	MimeType     string `json:"mime_type"`
	FileSize     int    `json:"file_size"`
}

type SendDocumentResponse struct {
	Ok     bool    `json:"ok"`
	Result Message `json:"result"`
}

type SendMessageResponse struct {
	Ok     bool    `json:"ok"`
	Result Message `json:"result"`
}
