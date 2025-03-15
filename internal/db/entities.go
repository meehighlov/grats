package db

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/meehighlov/grats/internal/config"
)

type BaseFields struct {
	ID        string
	CreatedAt string
	UpdatedAt string
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

func NewBaseFields() BaseFields {
	id := uuid.New().String()
	location, err := time.LoadLocation(config.Cfg().Timezone)
	if err != nil {
		slog.Error("error loading location by timezone, using system timezone, error: " + err.Error() + " NewEntityId: " + id)
	}
	now := time.Now().In(location).Format("02.01.2006T15:04:05")
	return BaseFields{id, now, now}
}

type User struct {
	// telegram user -> bot's user

	BaseFields

	TGId       string // id will be taken from telegram
	Name       string
	TGusername string
	ChatId     string // chatId - id of chat with user, bot uses it to send notification
	Birthday   string
	IsAdmin    int

	Friends []Friend
}

func (user *User) HasAdminAccess() bool {
	return user.IsAdmin == 1
}

func (user *User) GetTGUserName() string {
	if !strings.HasPrefix("@", user.TGusername) {
		return "@" + user.TGusername
	}
	return user.TGusername
}

func (user *User) FriendsListAsString() string {
	result := ""
	for _, friend := range user.Friends {
		result += friend.Name + " " + friend.BirthDay + "\n"
	}
	return result
}

type Friend struct {
	BaseFields

	// todo store timezone in friend table or somewere in db - for user's specific timezone

	Name           string
	UserId         string
	BirthDay       string
	ChatId         string
	notifyAt       string
	FilterNotifyAt string // this param is only for filtering
}

type Chat struct {
	BaseFields

	// more: https://core.telegram.org/bots/api#chat
	// may be one of: private, group, supergroup, channel
	// lowercase!
	// todo enum
	ChatType string

	BotInvitedBy        string
	ChatId              string
	GreetingTemplate    string

	// 0 off, 1 on
	SilentNotifications int
}

func (chat *Chat) IsAlreadySilent() bool {
	return chat.SilentNotifications == 1
}

func (chat *Chat) EnableSoundNotifications() {
	chat.SilentNotifications = 0
}

func (chat *Chat) DisableSoundNotifications() {
	chat.SilentNotifications = 1
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
	return &friend.notifyAt
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

type Access struct {
	BaseFields

	TGusername string
}

func (access *Access) GetTGUserName() string {
	if !strings.HasPrefix("@", access.TGusername) {
		return "@" + access.TGusername
	}
	return access.TGusername
}
