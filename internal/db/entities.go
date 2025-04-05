package db

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/meehighlov/grats/internal/common"
	"github.com/meehighlov/grats/internal/config"
)

const (
	SHORT_ID_LENGTH = 6
)

type BaseFields struct {
	ID        string `gorm:"primaryKey;type:string;column:id"`
	CreatedAt string `gorm:"not null;column:created_at;type:varchar"`
	UpdatedAt string `gorm:"not null;column:updated_at;type:varchar"`
}

func (b *BaseFields) RefresTimestamps() (created string, updated string, _ error) {
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " entityId: " + b.ID)
	}
	now := time.Now().In(location).Format("02.01.2006T15:04:05")
	if b.CreatedAt == "" {
		b.CreatedAt = now
	}
	b.UpdatedAt = now

	return b.CreatedAt, b.UpdatedAt, nil
}

func NewBaseFields(shortId bool) BaseFields {
	id := uuid.New().String()
	if shortId {
		id = GenerateShortID(SHORT_ID_LENGTH)
	}
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " NewEntityId: " + id)
	}
	now := time.Now().In(location).Format("02.01.2006T15:04:05")
	return BaseFields{id, now, now}
}

func NewEntity(table string) common.PaginatedEntity {
	if table == "user" {
		return &User{}
	}
	if table == "friend" {
		return &Friend{}
	}
	if table == "chat" {
		return &Chat{}
	}
	if table == "wish" {
		return &Wish{}
	}
	if table == "wish_list" {
		return &WishList{}
	}
	return nil
}

type User struct {
	// telegram user -> bot's user
	BaseFields

	TgId       string `gorm:"uniqueIndex;not null;column:tg_id;type:varchar"` // id will be taken from telegram
	Name       string `gorm:"not null;column:name"`
	TgUsername string `gorm:"not null;column:tg_username"`
	ChatId     string `gorm:"column:chat_id;type:varchar"` // chatId - id of chat with user, bot uses it to send notification
	Birthday   string `gorm:"column:birthday;type:varchar"`
	IsAdmin    bool   `gorm:"column:is_admin;type:boolean"`

	Friends []Friend `gorm:"foreignKey:UserId;references:ID"`
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

type Friend struct {
	BaseFields

	// todo store timezone in friend table or somewere in db - for user's specific timezone

	Name     string `gorm:"not null;type:varchar"`
	UserId   string `gorm:"not null;index;type:varchar"`
	BirthDay string `gorm:"column:birthday;type:varchar"`
	ChatId   string `gorm:"column:chat_id;type:varchar"`
	NotifyAt string `gorm:"column:notify_at;type:varchar"`

	FilterNotifyAt string `gorm:"-"`

	User User `gorm:"foreignKey:UserId;references:ID"`
}

func (Friend) TableName() string {
	return "friend"
}

func (friend *Friend) GetUserId() string {
	return friend.UserId
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

func (friend *Friend) BirthDayAsObj(format string) (time.Time, error) {
	parts := strings.Split(friend.BirthDay, ".")
	birtday_wo_year := strings.Join(parts[:2], ".")

	return time.Parse(format, birtday_wo_year)
}

func (friend *Friend) GetZodiacSign() (emoji, text string) {
	format := "02.01" // day.month
	birthday, err := friend.BirthDayAsObj(format)

	if err != nil {
		slog.Error("define zodiac sign error: cannot parse birthday: " + err.Error())
		return "🌙", "это луна"
	}

	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone during zodiac sign defenition, error: " + err.Error() + " friendId: " + friend.ID)
		return "🌙", "это луна"
	}

	border := func(month time.Month, day int) time.Time {
		return time.Date(0, month, day, 0, 0, 0, 0, location)
	}

	if birthday.After(border(time.March, 21)) && birthday.Before(border(time.April, 20)) {
		return "♈️", "овен"
	}
	if birthday.After(border(time.April, 20)) && birthday.Before(border(time.May, 21)) {
		return "♉", "телец"
	}
	if birthday.After(border(time.May, 21)) && birthday.Before(border(time.June, 22)) {
		return "♊", "близнецы"
	}
	if birthday.After(border(time.June, 22)) && birthday.Before(border(time.July, 23)) {
		return "♋", "рак"
	}
	if birthday.After(border(time.July, 23)) && birthday.Before(border(time.August, 23)) {
		return "♌", "лев"
	}
	if birthday.After(border(time.August, 23)) && birthday.Before(border(time.September, 23)) {
		return "♍", "дева"
	}
	if birthday.After(border(time.September, 23)) && birthday.Before(border(time.October, 24)) {
		return "♎", "весы"
	}
	if birthday.After(border(time.October, 24)) && birthday.Before(border(time.November, 22)) {
		return "♏", "скорпион"
	}
	if birthday.After(border(time.November, 22)) && birthday.Before(border(time.December, 22)) {
		return "♐", "стрелец"
	}
	if birthday.After(border(time.December, 22)) && birthday.Before(border(time.December, 31)) {
		return "♑", "козерог"
	}
	if birthday.Equal(border(time.December, 31)) {
		return "♑", "козерог"
	}
	if birthday.After(border(time.January, 1)) && birthday.Before(border(time.January, 20)) {
		return "♑", "козерог"
	}
	if birthday.After(border(time.January, 20)) && birthday.Before(border(time.February, 19)) {
		return "♒", "водолей"
	}
	if birthday.After(border(time.February, 19)) && birthday.Before(border(time.March, 21)) {
		return "♓", "рыбы"
	}

	slog.Error("zodiac sign was not defined by birthday: " + friend.BirthDay)
	return "🌙", "это луна"
}

func (friend *Friend) CountDaysToBirthday() int {
	// todo store timezone in friend table or somewere in db - for user's specific timezone
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}
	now := time.Now().In(location)
	notify, err := time.Parse("02.01.2006", *friend.GetNotifyAt())
	if err != nil {
		slog.Error("error parsing notify during count days to birthday: " + err.Error())
		return -1
	}

	diff := now.Sub(notify)
	diff_days := diff.Hours() / 24

	return int(diff_days)
}

func (friend *Friend) IsThisMonthAfterToday() bool {
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}
	now := strings.Split(time.Now().In(location).Format("02.01.2006"), ".")
	thisMonth := strings.Split(friend.BirthDay, ".")[1] == now[1]
	afterToday := strings.Split(friend.BirthDay, ".")[0] > now[0]

	return thisMonth && afterToday
}

func (friend *Friend) IsTodayBirthday() bool {
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}
	now := strings.Split(time.Now().In(location).Format("02.01.2006"), ".")
	bd := strings.Split(friend.BirthDay, ".")

	return now[0] == bd[0] && now[1] == bd[1]
}

func (friend *Friend) GetNotifyAt() *string {
	return &friend.NotifyAt
}

func (friend *Friend) RenewNotifayAt() (string, error) {
	format := "02.01" // day.month

	birthday, err := friend.BirthDayAsObj(format)

	if err != nil {
		slog.Error("notify date creation: cannot parse birthday: " + err.Error())
		return "", nil
	}

	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}

	today, err := time.Parse(format, time.Now().In(location).Format(format))

	if err != nil {
		slog.Error("notify date creation: cannot parse today date:" + err.Error())
		return "", nil
	}

	year := time.Now().In(location).Year()
	if today.After(birthday) || today.Equal(birthday) {
		year += 1
	}

	*friend.GetNotifyAt() = fmt.Sprintf(birthday.Format(format)+".%d", year)

	return *friend.GetNotifyAt(), nil
}

func (friend *Friend) UpdateNotifyAt() (string, error) {
	format := "02.01.2006" // day.month.year
	notifyAt, err := time.Parse(format, *friend.GetNotifyAt())

	if err != nil {
		slog.Error("error updating notfiy date: " + err.Error())
		return *friend.GetNotifyAt(), err
	}

	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " friendId: " + friend.ID)
	}

	// in case of duplicated call
	if notifyAt.Year() != time.Now().In(location).Year() {
		return *friend.GetNotifyAt(), nil
	}

	*friend.GetNotifyAt() = notifyAt.AddDate(1, 0, 0).Format(format)
	return *friend.GetNotifyAt(), nil
}

type Chat struct {
	BaseFields

	// more: https://core.telegram.org/bots/api#chat
	// may be one of: private, group, supergroup, channel
	// lowercase!
	// todo enum
	ChatType string `gorm:"column:chat_type;type:varchar"`

	BotInvitedById   string `gorm:"column:bot_invited_by_id;type:varchar"`
	ChatId           string `gorm:"uniqueIndex;not null;column:chat_id"`
	GreetingTemplate string `gorm:"column:greeting_template;type:varchar"`

	SilentNotifications bool `gorm:"column:silent_notifications;type:boolean"`
}

func (Chat) TableName() string {
	return "chat"
}

func (chat *Chat) GetUserId() string {
	return chat.BotInvitedById
}

func (chat *Chat) IsAlreadySilent() bool {
	return chat.SilentNotifications
}

func (chat *Chat) GetSilent() bool {
	return chat.SilentNotifications
}

func (chat *Chat) EnableSoundNotifications() {
	chat.SilentNotifications = false
}

func (chat *Chat) DisableSoundNotifications() {
	chat.SilentNotifications = true
}
