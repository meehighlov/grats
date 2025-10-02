package entities

import (
	"strings"
)

type User struct {
	// telegram user -> bot's user
	BaseFields

	TgId       string `gorm:"uniqueIndex;not null;column:tg_id;type:varchar"` // id will be taken from telegram
	Name       string `gorm:"not null;column:name;type:varchar"`
	TgUsername string `gorm:"not null;column:tg_username;type:varchar"`
	ChatId     string `gorm:"column:chat_id;type:varchar"` // chatId - id of chat with user, bot uses it to send notification
	IsAdmin    bool   `gorm:"column:is_admin;type:boolean"`
}

func (User) TableName() string {
	return "user"
}

func (user *User) GetUserId() string {
	return user.ID
}

func (user *User) HasAdminAccess() bool {
	return user.IsAdmin
}

func (user *User) GetTGUserName() string {
	if !strings.HasPrefix("@", user.TgUsername) {
		return "@" + user.TgUsername
	}
	return user.TgUsername
}

func (user *User) GetId() string {
	return user.ID
}

func (user *User) ButtonText() string {
	return user.Name
}
