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

func (wish *Wish) IsOzon() bool {
	ozonPrefix1 := "https://ozon.ru/"
	ozonPrefix2 := "https://www.ozon.ru/"
	return strings.HasPrefix(wish.Link, ozonPrefix1) || strings.HasPrefix(wish.Link, ozonPrefix2)
}

func (wish *Wish) IsWildberries() bool {
	wbPrefix1 := "https://wildberries.ru/"
	wbPrefix2 := "https://www.wildberries.ru/"
	return strings.HasPrefix(wish.Link, wbPrefix1) || strings.HasPrefix(wish.Link, wbPrefix2)
}

func (wish *Wish) IsYandexMarket() bool {
	yandexMarketPrefix1 := "https://market.yandex.ru/"
	yandexMarketPrefix2 := "https://www.market.yandex.ru/"
	return strings.HasPrefix(wish.Link, yandexMarketPrefix1) || strings.HasPrefix(wish.Link, yandexMarketPrefix2)
}

func (wish *Wish) IsAvito() bool {
	avitoPrefix1 := "https://avito.ru/"
	avitoPrefix2 := "https://www.avito.ru/"
	return strings.HasPrefix(wish.Link, avitoPrefix1) || strings.HasPrefix(wish.Link, avitoPrefix2)
}

func (wish *Wish) GetMarketplace() string {
	switch {
	case wish.IsOzon():
		return "ozon"
	case wish.IsWildberries():
		return "wildberries"
	case wish.IsYandexMarket():
		return "yandex market"
	case wish.IsAvito():
		return "avito"
	}
	return "смотреть на сайте"
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
