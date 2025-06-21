package db

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/meehighlov/grats/internal/common"
)

type BaseFields struct {
	ID        string `gorm:"primaryKey;column:id;type:varchar(36);default:uuid_generate_v4()"`
	CreatedAt string `gorm:"column:created_at;type:timestamp;default:now()"`
	UpdatedAt string `gorm:"column:updated_at;type:timestamp;default:now()"`
}

func (b *BaseFields) RefresTimestamps() (string, string, string) {
	if b.ID == "" {
		b.ID = GenerateShortID(6)
	}
	b.CreatedAt = GetCurrentTimestamp()
	b.UpdatedAt = GetCurrentTimestamp()

	return b.ID, b.CreatedAt, b.UpdatedAt
}

func NewBaseFields(withoutId bool) BaseFields {
	var id string
	if !withoutId {
		id = GenerateShortID(6)
	}

	return BaseFields{
		ID:        id,
		CreatedAt: GetCurrentTimestamp(),
		UpdatedAt: GetCurrentTimestamp(),
	}
}

func CreateEntity(table string) common.PaginatedEntity {
	switch table {
	case "wish":
		return &Wish{}
	case "wish_list":
		return &WishList{}
	case "user":
		return &User{}
	}
	slog.Error("CreateEntity", "unknown table", table)
	return nil
}

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

type WishList struct {
	BaseFields

	Name   string `gorm:"column:name;type:varchar"`
	UserId string `gorm:"not null;index;column:user_id;type:varchar"`
	ChatId string `gorm:"column:chat_id;type:varchar"`

	User User `gorm:"foreignKey:UserId;references:ID"`
}

func (WishList) TableName() string {
	return "wish_list"
}

func (wishList *WishList) GetUserId() string {
	return wishList.UserId
}

type Wish struct {
	BaseFields

	Name         string `gorm:"column:name;type:varchar"`
	ChatId       string `gorm:"column:chat_id;type:varchar"`
	UserId       string `gorm:"not null;index;column:user_id;type:varchar"`
	Link         string `gorm:"column:link;type:varchar"`
	ExecutorId   string `gorm:"column:executor_id;type:varchar"`
	Price        string `gorm:"column:price;type:varchar"`
	WishListId   string `gorm:"column:wish_list_id;type:varchar"`
	WishListName string `gorm:"-"`

	User     User     `gorm:"foreignKey:UserId;references:ID"`
	WishList WishList `gorm:"foreignKey:WishListId;references:ID"`
}

func (Wish) TableName() string {
	return "wish"
}

func (wish *Wish) GetUserId() string {
	return wish.UserId
}

func (wish *Wish) Info(executorId string) string {
	price := ""
	if wish.Price != "" {
		price = fmt.Sprintf(" - %s(RUB)", wish.Price)
	}
	msgLines := []string{
		fmt.Sprintf("✨ %s%s", wish.Name, price),
	}
	if wish.ExecutorId == "" {
		msgLines = append(msgLines, "🟢 желание пока не выбрали")
	} else {
		if wish.ExecutorId == executorId {
			msgLines = append(msgLines, "🎁 вы забронировали это желание")
		} else {
			msgLines = append(msgLines, "🎁 кто-то забронировал это желание")
		}
	}
	return strings.Join(msgLines, "\n\n")
}

func (wish *Wish) GetMarketplace() string {
	return "смотреть на сайте"
}

type Chat struct {
	BaseFields

	// more: https://core.telegram.org/bots/api#chat
	// may be one of: private, group, supergroup, channel
	// lowercase!
	// todo enum
	ChatType string `gorm:"column:chat_type;type:varchar"`

	BotInvitedById string `gorm:"column:bot_invited_by_id;type:varchar"`
	ChatId         string `gorm:"uniqueIndex;not null;column:chat_id;type:varchar"`
}

func (Chat) TableName() string {
	return "chat"
}

func (chat *Chat) GetUserId() string {
	return chat.BotInvitedById
}
