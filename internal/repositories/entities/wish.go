package entities

import (
	"fmt"
	"strings"
)

type Wish struct {
	BaseFields

	ChatId     string `gorm:"column:chat_id;type:varchar"`
	UserId     string `gorm:"not null;index;column:user_id;type:varchar"`
	WishListId string `gorm:"column:wish_list_id;type:varchar"`
	Link       string `gorm:"column:link;type:varchar"`
	Price      string `gorm:"column:price;type:varchar"`
	Name       string `gorm:"column:name;type:varchar"`
	ExecutorId string `gorm:"column:executor_id;type:varchar"`

	User     User     `gorm:"foreignKey:UserId;references:ID"`
	WishList WishList `gorm:"foreignKey:WishListId;references:ID"`
}

func (Wish) TableName() string {
	return "wish"
}

func (wish *Wish) GetUserId() string {
	return wish.UserId
}

func (wish *Wish) GetMarketplace(getSiteName func(string) (string, error)) string {
	name, err := getSiteName(wish.Link)
	if err != nil {
		return "смотреть на сайте"
	}
	return name
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

func (w *Wish) ButtonText() string {
	text := w.Name
	if w.Price != "" {
		text = w.Price + "(RUB)" + " " + w.Name
	}
	if w.ExecutorId != "" {
		text = "🔒 " + text
	}
	return text
}
